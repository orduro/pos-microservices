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
	ErrTenantNotFound       = errors.New("tenant not found in context")
	ErrItemNotFound         = errors.New("item not found")
	ErrItemAlreadyArchived  = errors.New("item already archived")
	ErrItemNotArchived      = errors.New("item not archived")
	ErrItemDuplicate        = errors.New("item with this name already exists in venue")
)
