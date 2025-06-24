package customer

import (
	"errors"

	"github.com/devsrivatsa/tavernDDD/domain"
	"github.com/google/uuid"
)

var (
	ErrInvalidPerson = errors.New("a customer must have a valid name")
)

type Customer struct {
	//person is the root entity of the customer aggregate
	person       *domain.Person
	products     []*domain.Item
	transactions []domain.Transaction
}

func NewCustomer(name string) (Customer, error) {
	if name == "" {
		return Customer{}, ErrInvalidPerson
	}
	person := &domain.Person{
		Name: name,
		ID:   uuid.New(),
	}

	return Customer{
		person:       person,
		products:     make([]*domain.Item, 0),
		transactions: make([]domain.Transaction, 0),
	}, nil
}

func (c *Customer) GetID() uuid.UUID {
	return c.person.ID
}

func (c *Customer) SetID(id uuid.UUID) {
	if c.person == nil {
		c.person = &domain.Person{}
	}
	c.person.ID = id
}

func (c *Customer) GetName() string {
	return c.person.Name
}

func (c *Customer) SetName(name string) {
	if c.person == nil {
		c.person = &domain.Person{}
	}
}
