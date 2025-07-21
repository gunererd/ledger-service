package mongo

import (
	"context"
	"ledger-service/internal/core/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransactionRepository struct {
	collection *mongo.Collection
}

func NewTransactionRepository(client *mongo.Client, dbName, collectionName string) *TransactionRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &TransactionRepository{
		collection: collection,
	}
}

func (r *TransactionRepository) Save(ctx context.Context, t types.Transaction) (string, error) {

	if t.Id == "" {
		t.Id = primitive.NewObjectID().Hex()
	}

	_, err := r.collection.InsertOne(ctx, t)
	if err != nil {
		return "", err
	}

	return t.Id, nil
}

func (r *TransactionRepository) GetOne(ctx context.Context, id string) (types.Transaction, error) {
	var transaction types.Transaction
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&transaction)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return types.Transaction{}, nil // Not found is OK
		}
		return types.Transaction{}, err // Actual error
	}
	return transaction, nil
}

func (r *TransactionRepository) GetManyForCustomer(ctx context.Context, id string) ([]types.Transaction, error) {
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"customer.id": id}, opts)
	if err != nil {
		return []types.Transaction{}, err
	}
	defer cursor.Close(ctx)

	results := []types.Transaction{}
	for cursor.Next(ctx) {
		var transaction types.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			return results, err // Return partial results with error
		}
		results = append(results, transaction)
	}
	return results, cursor.Err()
}

func (r *TransactionRepository) GetManyForRestaurant(ctx context.Context, id string) ([]types.Transaction, error) {
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})
	cursor, err := r.collection.Find(ctx, bson.M{"restaurant.id": id}, opts)
	if err != nil {
		return []types.Transaction{}, err
	}
	defer cursor.Close(ctx)

	results := []types.Transaction{}
	for cursor.Next(ctx) {
		var transaction types.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			return results, err // Return partial results with error
		}
		results = append(results, transaction)
	}
	return results, cursor.Err()
}
