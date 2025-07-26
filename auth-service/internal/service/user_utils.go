package service

import (
	"context"
	"fmt"
	"log"

	"github.com/orduro/pos-microservices/auth-service/internal/constants"
	"golang.org/x/crypto/bcrypt"
)

// this function generates and stores a token for a user
// then sends them an email with a link to the frontend
// that includes their token in the url query string
//
// can currently be used for sending password reset and email verification
// allowed modes:
// - constants.TokenTypeEmailVerification:
// - constants.TokenTypePasswordReset
func (s *UserService) sendTokenCallbackEmailToUser(email string, userID int64, mode string) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.ServiceCallTimeout)
	defer cancel()

	token, err := s.tokenService.CreateVerificationToken(ctx, userID, mode)
	if err != nil {
		log.Printf("failed to create verification token for user %d: %v", userID, err)
		return
	}

	// check mode and determine callback link
	callbackLink := ""
	switch mode {
	case constants.TokenTypeEmailVerification:
		callbackLink = fmt.Sprintf("%s/verification?token=%s", s.frontendURL, token)
	case constants.TokenTypePasswordReset:
		callbackLink = fmt.Sprintf("%s/password-reset?token=%s", s.frontendURL, token)
	default:
		log.Printf("token callback email can't be sent, invalid mode: %v", mode)
		return
	}

	if err := s.httpClient.SendCallbackEmail(ctx, email, callbackLink); err != nil {
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
