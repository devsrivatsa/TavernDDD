package memory

import (
	"errors"
	"testing"

	"github.com/devsrivatsa/tavernDDD/domain/customer"
	"github.com/google/uuid"
)

func TestMemoryStore_Get(t *testing.T) {
	type testCase struct {
		name          string
		id            uuid.UUID
		expectedError error
	}

	NewCustomer, err := customer.NewCustomer("Percy")
	if err != nil {
		t.Fatal(err)
	}
	id := NewCustomer.GetID()
	repo := MemoryStore{
		customers: map[uuid.UUID]customer.Customer{
			id: NewCustomer,
		},
	}
	testCases := []testCase{
		{
			name:          "no customer by id",
			id:            uuid.MustParse("f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			expectedError: customer.ErrCustomerNotFound,
		},
		{
			name:          "customer by id",
			id:            id,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := repo.Get(tc.id)
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}
		})
	}
}
