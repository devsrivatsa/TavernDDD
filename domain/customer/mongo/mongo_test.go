package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/devsrivatsa/tavernDDD/domain/customer"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	mongoConnectionString = "mongodb://localhost:27017"
	testTimeout           = 30 * time.Second
)

func setupTestRepo(t *testing.T) *MongoRepository {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	repo, err := New(ctx, mongoConnectionString)
	require.NoError(t, err, "Failed to create MongoDB repository")

	// Clean up any existing test data
	cleanupTestData(t, repo)

	return repo
}

func cleanupTestData(t *testing.T, repo *MongoRepository) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Drop the test collection to ensure clean state
	err := repo.customer.Drop(ctx)
	if err != nil {
		t.Logf("Warning: Failed to drop collection during cleanup: %v", err)
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name             string
		connectionString string
		expectError      bool
	}{
		{
			name:             "Valid connection string",
			connectionString: mongoConnectionString,
			expectError:      false,
		},
		{
			name:             "Invalid connection string",
			connectionString: "invalid://connection",
			expectError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
			defer cancel()

			repo, err := New(ctx, tt.connectionString)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, repo)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, repo)
				assert.NotNil(t, repo.db)
				assert.NotNil(t, repo.customer)
			}
		})
	}
}

func TestMongoRepository_Add(t *testing.T) {
	repo := setupTestRepo(t)
	defer cleanupTestData(t, repo)

	tests := []struct {
		name         string
		customerName string
		expectError  bool
	}{
		{
			name:         "Valid customer",
			customerName: "John Doe",
			expectError:  false,
		},
		{
			name:         "Another valid customer",
			customerName: "Jane Smith",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customer, err := customer.NewCustomer(tt.customerName)
			require.NoError(t, err)

			err = repo.Add(customer)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify the customer was actually added
				retrievedCustomer, err := repo.Get(customer.GetID())
				assert.NoError(t, err)
				assert.Equal(t, customer.GetID(), retrievedCustomer.GetID())
				assert.Equal(t, customer.GetName(), retrievedCustomer.GetName())
			}
		})
	}
}

func TestMongoRepository_Get(t *testing.T) {
	repo := setupTestRepo(t)
	defer cleanupTestData(t, repo)

	// Add a test customer first
	testCustomer, err := customer.NewCustomer("Test Customer")
	require.NoError(t, err)

	err = repo.Add(testCustomer)
	require.NoError(t, err)

	tests := []struct {
		name        string
		customerID  uuid.UUID
		expectError bool
	}{
		{
			name:        "Existing customer",
			customerID:  testCustomer.GetID(),
			expectError: false,
		},
		{
			name:        "Non-existing customer",
			customerID:  uuid.New(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retrievedCustomer, err := repo.Get(tt.customerID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, customer.Customer{}, retrievedCustomer)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.customerID, retrievedCustomer.GetID())
				assert.Equal(t, "Test Customer", retrievedCustomer.GetName())
			}
		})
	}
}

func TestMongoRepository_Update(t *testing.T) {
	repo := setupTestRepo(t)
	defer cleanupTestData(t, repo)

	// Add a test customer first
	originalCustomer, err := customer.NewCustomer("Original Name")
	require.NoError(t, err)

	err = repo.Add(originalCustomer)
	require.NoError(t, err)

	tests := []struct {
		name          string
		setupCustomer func() customer.Customer
		expectError   bool
	}{
		{
			name: "Update existing customer",
			setupCustomer: func() customer.Customer {
				updatedCustomer := originalCustomer
				updatedCustomer.SetName("Updated Name")
				return updatedCustomer
			},
			expectError: false,
		},
		{
			name: "Update non-existing customer",
			setupCustomer: func() customer.Customer {
				nonExistingCustomer, _ := customer.NewCustomer("Non Existing")
				return nonExistingCustomer
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customerToUpdate := tt.setupCustomer()

			err := repo.Update(customerToUpdate)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify the update was successful
				retrievedCustomer, err := repo.Get(customerToUpdate.GetID())
				assert.NoError(t, err)
				assert.Equal(t, customerToUpdate.GetName(), retrievedCustomer.GetName())
			}
		})
	}
}

func TestMongoRepository_Delete(t *testing.T) {
	repo := setupTestRepo(t)
	defer cleanupTestData(t, repo)

	// Add a test customer first
	testCustomer, err := customer.NewCustomer("Customer To Delete")
	require.NoError(t, err)

	err = repo.Add(testCustomer)
	require.NoError(t, err)

	tests := []struct {
		name        string
		customerID  uuid.UUID
		expectError bool
	}{
		{
			name:        "Delete existing customer",
			customerID:  testCustomer.GetID(),
			expectError: false,
		},
		{
			name:        "Delete non-existing customer",
			customerID:  uuid.New(),
			expectError: false, // MongoDB DeleteOne doesn't error if document doesn't exist
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(tt.customerID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify the customer was actually deleted (only for the first test case)
				if tt.name == "Delete existing customer" {
					_, err := repo.Get(tt.customerID)
					assert.Error(t, err, "Customer should not exist after deletion")
				}
			}
		})
	}
}

func TestMongoCustomer_ToAggregate(t *testing.T) {
	testID := uuid.New()
	testName := "Test Customer"

	mongoCustomer := mongoCustomer{
		ID:   testID,
		Name: testName,
	}

	aggregate := mongoCustomer.ToAggregate()

	assert.Equal(t, testID, aggregate.GetID())
	assert.Equal(t, testName, aggregate.GetName())
}

func TestNewFromCustomer(t *testing.T) {
	testCustomer, err := customer.NewCustomer("Test Customer")
	require.NoError(t, err)

	mongoCustomer := NewFromCustomer(testCustomer)

	assert.Equal(t, testCustomer.GetID(), mongoCustomer.ID)
	assert.Equal(t, testCustomer.GetName(), mongoCustomer.Name)
}

// Integration test that tests the full CRUD cycle
func TestMongoRepository_FullCRUDCycle(t *testing.T) {
	repo := setupTestRepo(t)
	defer cleanupTestData(t, repo)

	// Create
	originalCustomer, err := customer.NewCustomer("CRUD Test Customer")
	require.NoError(t, err)

	err = repo.Add(originalCustomer)
	require.NoError(t, err)

	// Read
	retrievedCustomer, err := repo.Get(originalCustomer.GetID())
	require.NoError(t, err)
	assert.Equal(t, originalCustomer.GetID(), retrievedCustomer.GetID())
	assert.Equal(t, originalCustomer.GetName(), retrievedCustomer.GetName())

	// Update
	retrievedCustomer.SetName("Updated CRUD Customer")
	err = repo.Update(retrievedCustomer)
	require.NoError(t, err)

	// Verify Update
	updatedCustomer, err := repo.Get(retrievedCustomer.GetID())
	require.NoError(t, err)
	assert.Equal(t, "Updated CRUD Customer", updatedCustomer.GetName())

	// Delete
	err = repo.Delete(updatedCustomer.GetID())
	require.NoError(t, err)

	// Verify Delete
	_, err = repo.Get(updatedCustomer.GetID())
	assert.Error(t, err, "Customer should not exist after deletion")
}
