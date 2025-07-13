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
	ServiceNameAuth = "auth-service"
)

const (
	MsgRegistrationSuccess = "User registered successfully. Please check your email for verification instructions."
	MsgVerificationSent    = "if an account with this email exists, a verification email has been sent."
)

const (
	ErrMsgInvalidJSON         = "invalid json value, cannot be read"
	ErrMsgInvalidCredentials  = "invalid username or password"
	ErrMsgUserExists          = "account with this email already exists"
	ErrMsgUserNotFound        = "user not found"
	ErrMsgUserAlreadyVerified = "user is already verified"
	ErrMsgInternalError       = "unable to process request"
	ErrMsgInvalidEmail        = "invalid email format"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 128
	MaxEmailLength    = 255
)

const (
	StatusHealthy = "alive"
)
