package tavern

import (
	"fmt"
	"log"

	"github.com/devsrivatsa/tavernDDD/services/order"
	"github.com/google/uuid"
)

type TavernConfiguration func(t *Tavern) error

type Tavern struct {
	orderService *order.OrderService
	//billing service
	billingService interface{}
}

func NewTavern(configs ...TavernConfiguration) (*Tavern, error) {
	t := &Tavern{}

	for _, config := range configs {
		if err := config(t); err != nil {
			return nil, err
		}
	}
	return t, nil
}

func WithOrderService(os *order.OrderService) TavernConfiguration {
	return func(t *Tavern) error {
		t.orderService = os
		return nil
	}
}

//if you have a billing service, you can add it to the tavern

func (t *Tavern) Order(customerID uuid.UUID, products []uuid.UUID) error {
	price, err := t.orderService.CreateOrder(customerID, products)
	if err != nil {
		return fmt.Errorf("error creating order: %w", err)
	}

	log.Printf("\nBill the customer %s for the amount of %.2f\n", customerID, price)

	return nil
}
