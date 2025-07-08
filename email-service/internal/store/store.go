package store

type Store struct {
	Email EmailRepository
}

func NewStore(emailConfig EmailConfig) Store {
	return Store{
		Email: NewEmailStore(emailConfig),
	}
}

type EmailRepository interface {
	SendEmail(req EmailRequest) (*EmailResponse, error)
	SendVerificationEmail(to, username, verificationURL string) error
	SendPasswordResetEmail(to, username, resetURL string) error
	SendWelcomeEmail(to, username string) error
}
