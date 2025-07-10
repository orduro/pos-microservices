package httpclient

import (
	"context"
	"log"
)

type EmailVerificationRequest struct {
	Email           string `json:"email"`
	VerificationURL string `json:"verification_url"`
}

func (c *Client) SendVerificationEmail(ctx context.Context, email, verificationURL string) error {
	reqBody := EmailVerificationRequest{
		Email:           email,
		VerificationURL: verificationURL,
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
