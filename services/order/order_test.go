package order

import (
	"testing"

	"github.com/devsrivatsa/tavernDDD/domain/product"
	"github.com/google/uuid"
)

func init_products(t *testing.T) []product.Product {
	beer, err := product.NewProduct("Beer", "A refreshing beer", 1.99)
	if err != nil {
		t.Fatalf("Error initializing product: %v \nCannot proceed with tests", err)
	}
	peanuts, err := product.NewProduct("Peanuts", "A delicious snack", 0.99)
	if err != nil {
		t.Fatalf("Error initializing product: %v \nCannot proceed with tests", err)
	}
	wine, err := product.NewProduct("Wine", "A fine wine", 5.99)
	if err != nil {
		t.Fatalf("Error initializing product: %v \nCannot proceed with tests", err)
	}

	return []product.Product{beer, peanuts, wine}
}

func TestOrder_NewOrderService(t *testing.T) {
	products := init_products(t)
	or, err := NewOrderService(
		WithMemoryCustomerRepository(),
		WithMemoryProductRepository(products),
	)
	if err != nil {
		t.Errorf("Error creating order service: %v", err)
	}
	t.Log("Order service created")

	customerID, err := or.AddCustomer("John Doe")
	if err != nil {
		t.Errorf("Error creating customer: %v", err)
	}

	t.Log("Customer created and added to the order service")
	order := []uuid.UUID{products[0].GetID()}

	_, err = or.CreateOrder(customerID, order)
	if err != nil {
		t.Fatalf("Error creating order: %v", err)
	}
	t.Log("Order created")
}
