package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/devsrivatsa/tavernDDD/domain/customer"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	db       *mongo.Database
	customer *mongo.Collection
}

type mongoCustomer struct {
	ID   uuid.UUID `bson:"_id"`
	Name string    `bson:"name"`
}

func NewFromCustomer(c customer.Customer) mongoCustomer {
	return mongoCustomer{
		ID:   c.GetID(),
		Name: c.GetName(),
	}
}

func (m mongoCustomer) ToAggregate() customer.Customer {
	customer := customer.Customer{}
	customer.SetID(m.ID)
	customer.SetName(m.Name)

	return customer
}

func New(ctx context.Context, connString string) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))
	if err != nil {
		return nil, err
	}

	db := client.Database("tavern")
	customers := db.Collection("customers")

	return &MongoRepository{
		db:       db,
		customer: customers,
	}, nil
}

func (mr *MongoRepository) Get(id uuid.UUID) (customer.Customer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := mr.customer.FindOne(ctx, bson.M{"_id": id})
	var c mongoCustomer
	if err := result.Decode(&c); err != nil {
		return customer.Customer{}, err
	}

	return c.ToAggregate(), nil
}

func (mr *MongoRepository) Add(c customer.Customer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	internal := NewFromCustomer(c)
	_, err := mr.customer.InsertOne(ctx, internal)
	if err != nil {
		return err
	}

	return nil
}

func (mr *MongoRepository) Update(c customer.Customer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//check if customer exists
	_, err := mr.Get(c.GetID())
	if err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}
	internal := NewFromCustomer(c)
	_, err = mr.customer.UpdateOne(ctx, bson.M{"_id": c.GetID()}, bson.M{"$set": internal})
	if err != nil {
		return err
	}

	return nil
}

func (mr *MongoRepository) Delete(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := mr.customer.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
