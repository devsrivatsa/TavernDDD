package memory

import (
	"fmt"
	"sync"

	"github.com/devsrivatsa/tavernDDD/domain/product"
	"github.com/google/uuid"
)

type MemoryProductRepository struct {
	products map[uuid.UUID]product.Product
	sync.Mutex
}

func New() *MemoryProductRepository {
	return &MemoryProductRepository{
		products: make(map[uuid.UUID]product.Product),
	}
}

func (m *MemoryProductRepository) GetAll() ([]product.Product, error) {
	var products []product.Product

	for _, product := range m.products {
		products = append(products, product)
	}

	return products, nil
}

func (m *MemoryProductRepository) GetByID(id uuid.UUID) (product.Product, error) {
	if prd, ok := m.products[id]; ok {
		return prd, nil
	}

	return product.Product{}, product.ErrProductNotFound
}

func (m *MemoryProductRepository) Update(prd product.Product) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.products[prd.GetID()]; !ok {
		return product.ErrProductNotFound
	}
	m.products[prd.GetID()] = prd

	return nil
}

func (m *MemoryProductRepository) Delete(id uuid.UUID) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.products[id]; !ok {
		return product.ErrProductNotFound
	}
	delete(m.products, id)

	return nil
}

func (m *MemoryProductRepository) Add(prd product.Product) error {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.products[prd.GetID()]; ok {
		return fmt.Errorf("error adding product %v due to error: %w", prd.GetItem(), product.ErrProductAlreadyExists)
	}
	m.products[prd.GetID()] = prd

	return nil
}
