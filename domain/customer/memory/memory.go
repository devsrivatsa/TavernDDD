package memory

import (
	"fmt"
	"sync"

	"github.com/devsrivatsa/tavernDDD/domain/customer"
	"github.com/google/uuid"
)

type MemoryStore struct {
	customers map[uuid.UUID]customer.Customer
	sync.Mutex
}

func New() *MemoryStore {
	return &MemoryStore{
		customers: make(map[uuid.UUID]customer.Customer),
	}
}
func (ms *MemoryStore) Get(id uuid.UUID) (customer.Customer, error) {
	if customer, ok := ms.customers[id]; ok {
		return customer, nil
	}
	return customer.Customer{}, customer.ErrCustomerNotFound
}
func (ms *MemoryStore) Add(c customer.Customer) error {
	if ms.customers == nil {
		ms.Lock()
		ms.customers = make(map[uuid.UUID]customer.Customer)
		ms.Unlock()
	}
	if _, ok := ms.customers[c.GetID()]; ok {
		return fmt.Errorf("customer already exists %w", customer.ErrUpdateCustomer)
	}
	ms.Lock()
	ms.customers[c.GetID()] = c
	ms.Unlock()

	return nil
}
func (ms *MemoryStore) Update(c customer.Customer) error {
	if _, ok := ms.customers[c.GetID()]; !ok {
		return fmt.Errorf("customer does not exist %w", customer.ErrCustomerNotFound)
	}
	ms.Lock()
	ms.customers[c.GetID()] = c
	ms.Unlock()

	return nil
}
