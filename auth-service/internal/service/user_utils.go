package service

import (
	"context"
	"fmt"
	"log"

	"github.com/orduro/pos-microservices/auth-service/internal/constants"
	"golang.org/x/crypto/bcrypt"
)

func (s *UserService) sendVerificationEmailToUserAsync(email string, userID int64, tokenType string) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.ServiceCallTimeout)
	defer cancel()

	token, err := s.tokenService.CreateVerificationToken(ctx, userID, tokenType)
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
