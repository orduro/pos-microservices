package store

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var v = validator.New()

var (
	ErrVenueNotFound        = errors.New("venue not found")
	ErrVenueAlreadyArchived = errors.New("venue already archived")
	ErrVenueNotArchived     = errors.New("venue not archived")
	ErrVenueDuplicate       = errors.New("venue with this name and address already exists")
)
