package order

import (
	"log"

	"github.com/devsrivatsa/tavernDDD/domain/customer"
	custMem "github.com/devsrivatsa/tavernDDD/domain/customer/memory"
	"github.com/devsrivatsa/tavernDDD/domain/product"
	prdMem "github.com/devsrivatsa/tavernDDD/domain/product/memory"
	"github.com/google/uuid"
)

// OrderConfiguration is a function that configures the order service
// it takes a pointer to the order service and returns an error
// (service configuration generator pattern)
// OrderConfiguration is an alias or the function name to the defined function. This is a type.
type OrderConfiguration func(os *OrderService) error

type OrderService struct {
	customers customer.CustomerRepository
	products  product.ProductRepository
}

// factory function to create a new order service
func NewOrderService(cfgs ...OrderConfiguration) (*OrderService, error) {
	os := &OrderService{}
	// apply all configurations to the order service
	for _, cfg := range cfgs {
		err := cfg(os)
		if err != nil {
			return nil, err
		}
	}

	return os, nil
}

// applies a customer repository to the order service - this is an OrderConfiguration(the type) function
// the reason we are using this function type - OrderConfiguration as the return type is because we want to chain functions together
func WithCustomerRepository(cr customer.CustomerRepository) OrderConfiguration {
	return func(os *OrderService) error {
		os.customers = cr
		return nil
	}
}

func WithMemoryCustomerRepository() OrderConfiguration {

	cr := custMem.New()
	return WithCustomerRepository(cr)
}

//the reason we want to do the above is because we can initantiate a NewOrderService with:
/*
	os, err := NewOrderService(
		//WithMemoryCustomerRepository(),
		WithMongoCustomerRepository(),
		WithMemoryProductRepository(),
		WithLogging("debug"),
		WithTracing("123"),
	)

	The reason this is useful is that in future we can just change the implementation of the repository.
	E.g. if we want to use a sql database instead.
	Or we can just change the implementation of the logging or tracing to some other library or service.

*/

func WithMemoryProductRepository(products []product.Product) OrderConfiguration {

	return func(os *OrderService) error {
		pr := prdMem.New()
		for _, newPrd := range products {
			err := pr.Add(newPrd)
			if err != nil {
				return err
			}
		}
		os.products = pr

		return nil
	}
}

func (o *OrderService) CreateOrder(curstomerID uuid.UUID, productsID []uuid.UUID) (float64, error) {
	//fetch the customer

	customer, err := o.customers.Get(curstomerID)
	if err != nil {
		log.Printf("error fetching customer: %v", err)
		return 0, err
	}
	//fetch the products
	var products []product.Product
	var totalPrice float64
	for _, id := range productsID {
		prd, err := o.products.GetByID(id)
		if err != nil {
			log.Printf("error fetching product: %v", err)
			return 0, err
		}
		products = append(products, prd)
		totalPrice += prd.GetPrice()
	}
	log.Printf("Customer %s is ordering %d products for a total of %.2f", customer.GetName(), len(products), totalPrice)

	return totalPrice, nil
}

func (o *OrderService) AddCustomer(name string) (uuid.UUID, error) {
	c, err := customer.NewCustomer(name)
	if err != nil {
		return uuid.Nil, err
	}
	err = o.customers.Add(c)
	if err != nil {
		return uuid.Nil, err
	}

	return c.GetID(), nil
}
