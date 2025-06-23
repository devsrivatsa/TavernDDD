package product

import (
	"errors"

	"github.com/devsrivatsa/tavernDDD/entity"
	"github.com/google/uuid"
)

var (
	ErrMissingValue         = errors.New("missing important values")
	ErrProductNotFound      = errors.New("no such product found")
	ErrProductAlreadyExists = errors.New("product already exists")
)

type Product struct {
	item     *entity.Item
	price    float64
	quantity int
}

// factory function to create a new product
func NewProduct(name, description string, price float64) (Product, error) {
	if name == "" || description == "" {
		return Product{}, ErrMissingValue
	}
	return Product{
		item: &entity.Item{
			ID:          uuid.New(),
			Name:        name,
			Description: description,
		},
		price:    price,
		quantity: 1,
	}, nil
}

//these methods depend on what you need to expose

func (p Product) GetID() uuid.UUID {
	return p.item.ID
}
func (p Product) GetItem() *entity.Item {
	return p.item
}
func (p Product) GetPrice() float64 {
	return p.price
}
