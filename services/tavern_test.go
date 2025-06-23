package services

import (
	"testing"

	"github.com/devsrivatsa/tavernDDD/domain/customer"
	"github.com/google/uuid"
)

func TestTavern_Order(t *testing.T) {
	products := init_products(t)
	ordSrvc, err := NewOrderService(
		WithMemoryProductRepository(products),
		WithMemoryCustomerRepository(),
	)
	if err != nil {
		t.Fatalf("%v: Error creating order service: %v", t.Name(), err)
	}

	tavern, err := NewTavern(WithOrderService(ordSrvc))
	if err != nil {
		t.Fatalf("%v: Error creating tavern: %v", t.Name(), err)
	}

	c, err := customer.NewCustomer("ssriva")
	if err != nil {
		t.Fatalf("%v: Error creating customer: %v", t.Name(), err)
	}
	if err := ordSrvc.customers.AddCustomer(c); err != nil {
		t.Fatalf("%v: Error adding customer: %v", t.Name(), err)
	}
	order := []uuid.UUID{products[0].GetID()}
	err = tavern.Order(c.GetID(), order)
	if err != nil {
		t.Fatalf("%v: Error ordering: %v", t.Name(), err)
	}
	t.Logf("%v: Order successful", t.Name())
}
