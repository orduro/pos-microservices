package store

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

var v = validator.New()

var (
	ErrEmailRequired           = errors.New("email is required")
	ErrEmailInvalid            = errors.New("email format is invalid")
	ErrSubjectRequired         = errors.New("subject is required")
	ErrBodyRequired            = errors.New("body is required")
	ErrRecipientsRequired      = errors.New("at least one recipient is required")
	ErrEmailTypeInvalid        = errors.New("email type is invalid")
	ErrUsernameRequired        = errors.New("username is required")
	ErrVerificationURLRequired = errors.New("verification URL is required")
	ErrResetURLRequired        = errors.New("reset URL is required")
)

type EmailValidationRequest struct {
	To        []string  `validate:"required,min=1,dive,email"`
	Subject   string    `validate:"required,min=1,max=200"`
	Body      string    `validate:"required,min=1,max=10000"`
	EmailType EmailType `validate:"required"`
}

type VerificationEmailRequest struct {
	Email           string `validate:"required,email"`
	VerificationURL string `validate:"required,url"`
}

type PasswordResetEmailRequest struct {
	Email    string `validate:"required,email"`
	Username string `validate:"required,min=1,max=100"`
	ResetURL string `validate:"required,url"`
}

type WelcomeEmailRequest struct {
	Email    string `validate:"required,email"`
	Username string `validate:"required,min=1,max=100"`
}

func ValidateEmailRequest(req EmailRequest) error {
	validationReq := EmailValidationRequest{
		To:        req.To,
		Subject:   req.Subject,
		Body:      req.Body,
		EmailType: req.EmailType,
	}

	if err := v.Struct(validationReq); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, validationErr := range validationErrors {
				switch validationErr.Field() {
				case "To":
					if validationErr.Tag() == "required" || validationErr.Tag() == "min" {
						return ErrRecipientsRequired
					}
					if validationErr.Tag() == "email" {
						return ErrEmailInvalid
					}

				case "Subject":
					return ErrSubjectRequired

				case "Body":
					return ErrBodyRequired

				case "EmailType":
					return ErrEmailTypeInvalid
				}
			}
		}
		return err
	}

	if !IsValidEmailType(req.EmailType) {
		return ErrEmailTypeInvalid
	}

	return nil
}

func ValidateVerificationEmailRequest(email, verificationURL string) error {
	req := VerificationEmailRequest{
		Email:           email,
		VerificationURL: verificationURL,
	}

	if err := v.Struct(req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, validationErr := range validationErrors {
				switch validationErr.Field() {
				case "Email":
					return ErrEmailInvalid

				case "VerificationURL":
					return ErrVerificationURLRequired
				}
			}
		}
		return err
	}
	return nil
}

func ValidatePasswordResetEmailRequest(email, username, resetURL string) error {
	req := PasswordResetEmailRequest{
		Email:    email,
		Username: username,
		ResetURL: resetURL,
	}

	if err := v.Struct(req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, validationErr := range validationErrors {
				switch validationErr.Field() {
				case "Email":
					return ErrEmailInvalid

				case "Username":
					return ErrUsernameRequired

				case "ResetURL":
					return ErrResetURLRequired
				}
			}
		}
		return err
	}
	return nil
}

func ValidateWelcomeEmailRequest(email, username string) error {
	req := WelcomeEmailRequest{
		Email:    email,
		Username: username,
	}

	if err := v.Struct(req); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, validationErr := range validationErrors {
				switch validationErr.Field() {
				case "Email":
					return ErrEmailInvalid

				case "Username":
					return ErrUsernameRequired
				}
			}
		}
		return err
	}
	return nil
}

func SanitiseString(input string) string {
	cleaned := strings.ReplaceAll(input, "\x00", "")
	cleaned = strings.TrimSpace(cleaned)

	return cleaned
}

func IsValidEmailType(emailType EmailType) bool {
	switch emailType {
	case EmailTypeVerification, EmailTypePasswordReset, EmailTypeWelcome, EmailTypeNotification:
		return true

	default:
		return false
	}
}
