package customer_test

import (
	"errors"
	"testing"

	"github.com/devsrivatsa/tavernDDD/domain/customer"
)

func TestCustomer_NewCustomer(t *testing.T) {
	type testCase struct {
		test          string
		name          string
		expectedError error
	}

	testcases := []testCase{
		{
			test:          "Empty name validation",
			name:          "",
			expectedError: customer.ErrInvalidPerson,
		},
		{
			test:          "Valid name",
			name:          "John Doe",
			expectedError: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.test, func(t *testing.T) {
			_, err := customer.NewCustomer(tc.name)
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error %v, got %v", tc.expectedError, err)
			}
		})
	}
}
