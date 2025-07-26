package httpclient

import (
	"context"
	"log"
)

type EmailVerificationRequest struct {
	Email           string `json:"email"`
	VerificationURL string `json:"verification_url"`
}

func (c *Client) SendVerificationEmail(ctx context.Context, email, callbackURL string) error {
	reqBody := EmailVerificationRequest{
		Email:           email,
		VerificationURL: callbackURL,
	}

	request := ServiceRequest{
		Method: "POST",
		Path:   "/api/email/verification",
		Body:   reqBody,
	}

	resp, err := c.CallService(ctx, request)
	if err != nil {
		return err
	}

	log.Printf("response from communication service: %+v", resp)

	return nil
}

type PasswordResetEmailRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	ResetURL string `json:"reset_url"`
}

func (c *Client) SendPasswordResetEmail(ctx context.Context, email, callbackURL string) error {
	reqBody := PasswordResetEmailRequest{
		Email:    email,
		Username: email,
		ResetURL: callbackURL,
	}

	request := ServiceRequest{
		Method: "POST",
		Path:   "/api/email/password-reset",
		Body:   reqBody,
	}

	resp, err := c.CallService(ctx, request)
	if err != nil {
		return err
	}

	log.Printf("response from communication service: %+v", resp)

	return nil
}
