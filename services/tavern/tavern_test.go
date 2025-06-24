package tavern

import (
	"testing"

	"github.com/devsrivatsa/tavernDDD/domain/product"
	"github.com/devsrivatsa/tavernDDD/services/order"
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

func TestTavern_Order(t *testing.T) {
	products := init_products(t)
	ordSrvc, err := order.NewOrderService(
		order.WithMemoryProductRepository(products),
		order.WithMemoryCustomerRepository(),
	)
	if err != nil {
		t.Fatalf("%v: Error creating order service: %v", t.Name(), err)
	}

	tavern, err := NewTavern(WithOrderService(ordSrvc))
	if err != nil {
		t.Fatalf("%v: Error creating tavern: %v", t.Name(), err)
	}

	customerID, err := ordSrvc.AddCustomer("John Doe")
	if err != nil {
		t.Fatalf("%v: Error adding customer: %v", t.Name(), err)
	}
	order := []uuid.UUID{products[0].GetID()}
	err = tavern.Order(customerID, order)
	if err != nil {
		t.Fatalf("%v: Error ordering: %v", t.Name(), err)
	}
	t.Logf("%v: Order successful", t.Name())
}
