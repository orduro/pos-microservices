package store

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

// email sending request
type EmailRequest struct {
	To        []string  `json:"to"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	EmailType EmailType `json:"email_type"`
}

// response after sending an email
type EmailResponse struct {
	Success   bool      `json:"success"`
	MessageID string    `json:"message_id,omitempty"`
	Error     string    `json:"error,omitempty"`
	SentAt    time.Time `json:"sent_at"`
}

// different types of emails
type EmailType string

const (
	EmailTypeVerification  EmailType = "verification"
	EmailTypePasswordReset EmailType = "password_reset"
	EmailTypeWelcome       EmailType = "welcome"
	EmailTypeNotification  EmailType = "notification"
)

// email configuration
type EmailConfig struct {
	FromEmail    string
	FromName     string
	ReplyToEmail string
	Provider     string
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPass     string
}

// EmailStore implements EmailRepository
type EmailStore struct {
	config EmailConfig
	dialer *gomail.Dialer
}

func NewEmailStore(config EmailConfig) EmailRepository {
	store := &EmailStore{
		config: config,
	}

	// setup email dialer for gomail
	if config.Provider == "gomail" {
		store.dialer = gomail.NewDialer(
			config.SMTPHost,
			config.SMTPPort,
			config.SMTPUser,
			config.SMTPPass,
		)
	}

	return store
}

// send email using gomail
func (s *EmailStore) SendEmail(req EmailRequest) (*EmailResponse, error) {
	log.Printf("[Send Email] Starting to send %s email to: %v", req.EmailType, req.To)

	if err := ValidateEmailRequest(req); err != nil {
		log.Printf("[Send Email] Validation failed for %s email: %v", req.EmailType, err)

		return &EmailResponse{
			Success: false,
			Error:   err.Error(),
			SentAt:  time.Now(),
		}, err
	}

	req.Subject = SanitiseString(req.Subject)
	req.Body = SanitiseString(req.Body)

	if s.config.Provider == "simulate" {
		return s.simulateEmail(req)
	}

	log.Printf("[Send Email] Sending via SMTP (Provider: %s, Host: %s)", s.config.Provider, s.config.SMTPHost)

	// create email message
	m := gomail.NewMessage()

	// set headers
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail))
	m.SetHeader("To", req.To...)
	m.SetHeader("Subject", req.Subject)
	m.SetHeader("Reply-To", s.config.ReplyToEmail)

	// set body
	if s.isHTML(req.Body) {
		m.SetBody("text/html", req.Body)
	} else {
		m.SetBody("text/plain", req.Body)
	}

	startTime := time.Now()

	// send email
	if err := s.dialer.DialAndSend(m); err != nil {
		duration := time.Since(startTime)
		log.Printf("[Send Email] Failed to send %s email to %v after %v: %v", req.EmailType, req.To, duration, err)

		return &EmailResponse{
			Success: false,
			Error:   err.Error(),
			SentAt:  time.Now(),
		}, err
	}

	duration := time.Since(startTime)
	messageID := fmt.Sprintf("gomail_%d", time.Now().Unix())
	log.Printf("[Send Email] Successfully sent %s email to %v in %v (MessageID: %s)", req.EmailType, req.To, duration, messageID)

	return &EmailResponse{
		Success:   true,
		MessageID: messageID,
		SentAt:    time.Now(),
	}, nil
}

func (s *EmailStore) simulateEmail(req EmailRequest) (*EmailResponse, error) {
	log.Printf("   [SIMULATED] Email Details:")
	log.Printf("   From: %s <%s>", s.config.FromName, s.config.FromEmail)
	log.Printf("   To: %s", strings.Join(req.To, ", "))
	log.Printf("   Reply-To: %s", s.config.ReplyToEmail)
	log.Printf("   Subject: %s", req.Subject)
	log.Printf("   Type: %s", req.EmailType)
	log.Printf("   Body: %s", req.Body)

	return &EmailResponse{
		Success:   true,
		MessageID: fmt.Sprintf("sim_%d", time.Now().Unix()),
		SentAt:    time.Now(),
	}, nil
}

func (s *EmailStore) isHTML(content string) bool {
	return strings.Contains(content, "<") && strings.Contains(content, ">")
}

// send verification email
func (s *EmailStore) SendVerificationEmail(to, verificationURL string) error {
	log.Printf("[Verification] Preparing verification email for: %s", to)

	if err := ValidateVerificationEmailRequest(to, verificationURL); err != nil {
		log.Printf("[Verification] Validation failed for %s: %v", to, err)
		return err
	}

	verificationURL = SanitiseString(verificationURL)

	req := EmailRequest{
		To:        []string{to},
		Subject:   "Verify your email address",
		Body:      fmt.Sprintf(`<h1>Welcome!</h1><p>Please verify your email by clicking <a href="%s">here</a></p>`, verificationURL),
		EmailType: EmailTypeVerification,
	}

	_, err := s.SendEmail(req)
	if err != nil {
		log.Printf("[Verification] Failed to send verification email to %s: %v", to, err)
	} else {
		log.Printf("[Verification] Verification email queued successfully for %s", to)
	}
	return err
}

// send password reset email
func (s *EmailStore) SendPasswordResetEmail(to, username, resetURL string) error {
	log.Printf("[Password Reset] Preparing password reset email for: %s (username: %s)", to, username)

	if err := ValidatePasswordResetEmailRequest(to, username, resetURL); err != nil {
		log.Printf("[Password Reset] Validation failed for %s: %v", to, err)
		return err
	}

	username = SanitiseString(username)
	resetURL = SanitiseString(resetURL)

	req := EmailRequest{
		To:        []string{to},
		Subject:   "Reset your password",
		Body:      fmt.Sprintf(`<h1>Password Reset</h1><p>Hi %s, click <a href="%s">here</a> to reset your password</p>`, username, resetURL),
		EmailType: EmailTypePasswordReset,
	}

	_, err := s.SendEmail(req)
	if err != nil {
		log.Printf("[Password Reset] Failed to send password reset email to %s: %v", to, err)
	} else {
		log.Printf("[Password Reset] Password reset email queued successfully for %s", to)
	}
	return err
}

// send welcome email
func (s *EmailStore) SendWelcomeEmail(to, username string) error {
	log.Printf("[Welcome] Preparing welcome email for: %s (username: %s)", to, username)

	if err := ValidateWelcomeEmailRequest(to, username); err != nil {
		log.Printf("[Welcome] Validation failed for %s: %v", to, err)
		return err
	}

	username = SanitiseString(username)

	req := EmailRequest{
		To:        []string{to},
		Subject:   "Welcome to Orduro!",
		Body:      fmt.Sprintf(`<h1>Welcome %s!</h1><p>Thanks for joining. We're excited to have you!</p>`, username),
		EmailType: EmailTypeWelcome,
	}

	_, err := s.SendEmail(req)
	if err != nil {
		log.Printf("[Welcome] Failed to send welcome email to %s: %v", to, err)
	} else {
		log.Printf("[Welcome] Welcome email queued successfully for %s", to)
	}
	return err
}
