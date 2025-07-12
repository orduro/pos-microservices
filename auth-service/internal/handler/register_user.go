package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/orduro/common/json"
	"github.com/orduro/pos-microservices/auth-service/internal/store"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// receive user registration details
	registrationDetails := store.UserRegistrationDetails{}
	if err := json.Read(r, &registrationDetails); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, "unable to read registration details")
		log.Printf("unable to read registration details: %v", err)
		return
	}

	// validate registration details
	if err := registrationDetails.Validate(); err != nil {
		json.WriteError(w, r, http.StatusBadRequest, fmt.Sprintf("invalid username or password: %v", err))
		log.Printf("invalid username or password: %v", err)
		return
	}

	// check if user already exists
	existingUser, err := h.store.Users.GetByEmail(r.Context(), registrationDetails.Email)
	if err != nil && !errors.Is(err, store.ErrUserNotFound) {
		json.WriteError(w, r, http.StatusInternalServerError, "unable to check existing user")
		log.Printf("failed to check existing user: %v", err)
		return
	}
	if existingUser != nil {
		json.WriteError(w, r, http.StatusConflict, "user with this email already exists")
		return
	}

	hashedPassword, err := hashPassword(registrationDetails.Password)
	if err != nil {
		json.WriteError(w, r, http.StatusInternalServerError, "unable to process registration")
		log.Printf("failed to hash password: %v", err)
		return
	}

	// create and store new user
	user := store.User{
		Email:        registrationDetails.Email,
		PasswordHash: hashedPassword,
		FirstName:    "",
		LastName:     "",
		IsVerified:   false,
	}
	if err := h.store.Users.Create(r.Context(), &user); err != nil {
		json.WriteError(w, r, http.StatusInternalServerError, "unable to create user")
		log.Printf("failed to create user: %v", err)
		return
	}

	// generate verification token
	token, err := generateVerificationToken()
	if err != nil {
		json.WriteError(w, r, http.StatusInternalServerError, "unable to generate verification token")
		log.Printf("failed to generate verification token: %v", err)
		return
	}

	// store verification token in table
	// verification token lasts 24 hours
	verificationToken := store.VerificationToken{
		UserID:    user.ID,
		Token:     token,
		TokenType: "email_verification",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := h.store.VerificationTokens.Create(r.Context(), &verificationToken); err != nil {
		json.WriteError(w, r, http.StatusInternalServerError, "unable to create verification token")
		log.Printf("failed to create verification token: %v", err)
		return
	}

	// generate verification link and send verification email in background
	go func() {
		// create new context with timeout for this background operation
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// generate verification link to verify with current token
		verificationLink := fmt.Sprintf("%s/verification?token=%s", h.frontendAdminURL, token)

		err = h.httpClient.SendVerificationEmail(ctx, registrationDetails.Email, verificationLink)
		if err != nil {
			log.Printf("failed to send verification email: %v", err)
		}
	}()

	json.Write(w, http.StatusCreated, map[string]any{
		"message": "User registered successfully. Please check your email for verification instructions.",
		"email":   user.Email,
	})
}

// TODO: verify user (send email)
func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {}

// TODO: get JWT token
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {}

// TODO: validate JWT for other services
func (h *Handler) Validate(w http.ResponseWriter, r *http.Request) {}

// TODO: get current user info
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateVerificationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
