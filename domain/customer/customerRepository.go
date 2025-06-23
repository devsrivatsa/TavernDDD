package customer

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrCustomerNotFound    = errors.New("customer not found")
	ErrFailedToAddCustomer = errors.New("failed to add customer")
	ErrUpdateCustomer      = errors.New("failed to update customer")
)

type CustomerRepository interface {
	Get(id uuid.UUID) (Customer, error)
	Add(customer Customer) error
	Update(customer Customer) error
}
