package constants

import "time"

const (
	TokenTypeEmailVerification = "email_verification"
	TokenTypePasswordReset     = "password_reset"
)

const (
	TokenExpiryDuration = 24 * time.Hour
	TokenByteLength     = 32
)

const (
	ServiceCallTimeout = 30 * time.Second
	RequestTimeout     = 60 * time.Second
)

const (
	MsgRegistrationSuccess = "user registered successfully. please check your email for verification instructions."
	MsgVerificationSent    = "if an account with this email exists, a verification email has been sent."
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 128
	MaxEmailLength    = 255
)
