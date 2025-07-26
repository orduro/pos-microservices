package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/orduro/pos-microservices/auth-service/internal/constants"
	"github.com/orduro/pos-microservices/auth-service/internal/httpclient"
	"github.com/orduro/pos-microservices/auth-service/internal/store"
)

type UserService struct {
	user         store.UserRepository
	tokenService *TokenService
	jwtService   *JWTService
	httpClient   *httpclient.Client
	frontendURL  string
}

func NewUserService(
	userStore store.UserRepository,
	tokenService *TokenService,
	jwtService *JWTService,
	httpClient *httpclient.Client,
	frontendURL string,
) *UserService {
	return &UserService{
		user:         userStore,
		tokenService: tokenService,
		jwtService:   jwtService,
		httpClient:   httpClient,
		frontendURL:  frontendURL,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, details store.UserRegistrationDetails) (*store.User, error) {
	// check if user already exists
	existingUser, err := s.user.GetByEmail(ctx, details.Email)
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
	}

	if err := s.user.Create(ctx, &user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// send verification email asynchronously
	go s.sendVerificationEmail(details.Email, user.ID)

	return &user, nil
}

func (s *UserService) ResendVerificationEmail(ctx context.Context, email string) error {
	user, err := s.user.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			// don't reveal if user exists for security reasons
			log.Printf("account with email: %s doesn't exist, no emails will be sent.", email)
			return nil
		}
		return fmt.Errorf("failed to check user: %w", err)
	}

	if user.IsVerified {
		return errors.New("user is already verified")
	}

	// send verification email asynchronously
	go s.sendVerificationEmail(email, user.ID)

	return nil
}

func (s *UserService) VerifyUser(ctx context.Context, tokenStr string) error {
	token, err := s.tokenService.ValidateToken(ctx, tokenStr, constants.TokenTypeEmailVerification)
	if err != nil {
		return err
	}

	// check if user is already verified before attempting to mark as verified
	user, err := s.user.GetById(ctx, token.UserID)
	if err != nil {
		return err
	}

	if user.IsVerified {
		return errors.New("user is already verified")
	}

	if err := s.user.MarkAsVerified(ctx, token.UserID); err != nil {
		return err
	}

	return nil
}

func (s *UserService) LoginUser(ctx context.Context, details store.UserLoginDetails) (*TokenPair, error) {
	// authenticate user by checking email and password matches
	user, err := s.AuthenticateUser(ctx, details.Email, details.Password)
	if err != nil {
		return nil, err
	}

	// generate token pair
	tokenPair, err := s.jwtService.GenerateTokenPairWithTenant(user.ID, user.Email, &user.TenantID)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := s.user.UpdateLastLogin(context.Background(), user.ID); err != nil {
			log.Printf("failed to update last login: %v", err)
		}
	}()

	return tokenPair, nil

}

func (s *UserService) RefreshTokens(ctx context.Context, refreshToken string) (*TokenPair, error) {
	return s.jwtService.RefreshToken(refreshToken)
}

func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (*store.User, error) {
	user, err := s.user.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !s.verifyPassword(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsVerified {
		return nil, errors.New("email not verified")
	}

	return user, nil
}
