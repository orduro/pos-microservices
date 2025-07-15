package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/orduro/pos-microservices/auth-service/internal/constants"
	"github.com/orduro/pos-microservices/auth-service/internal/httpclient"
	"github.com/orduro/pos-microservices/auth-service/internal/store"
)

type UserService struct {
	userStore    store.UserRepository
	tokenService *TokenService
	httpClient   *httpclient.Client
	frontendURL  string
}

func NewUserService(
	userStore store.UserRepository,
	tokenService *TokenService,
	httpClient *httpclient.Client,
	frontendURL string,
) *UserService {
	return &UserService{
		userStore:    userStore,
		tokenService: tokenService,
		httpClient:   httpClient,
		frontendURL:  frontendURL,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, details store.UserRegistrationDetails) (*store.User, error) {
	// check if user already exists
	existingUser, err := s.userStore.GetByEmail(ctx, details.Email)
	if err != nil && !errors.Is(err, store.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, store.ErrUserExists
	}

	// hash password
	hashedPassword, err := s.hashPassword(details.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// create user
	user := store.User{
		Email:        details.Email,
		PasswordHash: hashedPassword,
		FirstName:    "",
		LastName:     "",
		IsVerified:   false,
	}

	if err := s.userStore.Create(ctx, &user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// send verification email asynchronously
	go s.sendVerificationEmailAsync(details.Email, user.ID)

	return &user, nil
}

func (s *UserService) ResendVerificationEmail(ctx context.Context, email string) error {
	user, err := s.userStore.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			// don't reveal if user exists for security reasons
			return nil
		}
		return fmt.Errorf("failed to check user: %w", err)
	}

	if user.IsVerified {
		return errors.New("user is already verified")
	}

	// send verification email asynchronously
	go s.sendVerificationEmailAsync(email, user.ID)

	return nil
}

func (s *UserService) VerifyUser(ctx context.Context, tokenStr string) error {
	token, err := s.tokenService.ValidateToken(ctx, tokenStr, constants.TokenTypeEmailVerification)
	if err != nil {
		return fmt.Errorf("invalid verification token: %w", err)
	}

	// check if user is already verified before attempting to mark as verified
	user, err := s.userStore.GetById(ctx, token.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user.IsVerified {
		return errors.New("user is already verified")
	}

	if err := s.userStore.MarkAsVerified(ctx, token.UserID); err != nil {
		return fmt.Errorf("failed to mark user as verified: %w", err)
	}

	return nil
}

func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (*store.User, error) {
	user, err := s.userStore.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !s.verifyPassword(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsVerified {
		return nil, errors.New("email not verified")
	}

	return user, nil
}

func (s *UserService) sendVerificationEmailAsync(email string, userID int64) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.ServiceCallTimeout)
	defer cancel()

	token, err := s.tokenService.CreateVerificationToken(ctx, userID)
	if err != nil {
		log.Printf("failed to create verification token for user %d: %v", userID, err)
		return
	}

	verificationLink := fmt.Sprintf("%s/verification?token=%s", s.frontendURL, token)

	if err := s.httpClient.SendVerificationEmail(ctx, email, verificationLink); err != nil {
		log.Printf("failed to send verification email to %s: %v", email, err)
	} else {
		log.Printf("verification email sent successfully to %s", email)
	}
}

func (s *UserService) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (s *UserService) verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
