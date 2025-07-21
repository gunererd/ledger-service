package mongo

import (
	"context"
	"ledger-service/internal/core/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BalanceRepository struct {
	collection *mongo.Collection
}

func NewBalanceRepository(client *mongo.Client, dbName, collectionName string) *BalanceRepository {
	collection := client.Database(dbName).Collection(collectionName)
	return &BalanceRepository{
		collection: collection,
	}
}

func (r *BalanceRepository) GetBalance(ctx context.Context, userId string) (types.Balance, error) {
	var balance types.Balance
	err := r.collection.FindOne(ctx, bson.M{"userid": userId}).Decode(&balance)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return types.Balance{UserId: userId, Amount: 0, TotalCommission: 0}, nil
		}
		return types.Balance{}, err
	}
	return balance, nil
}

func (r *BalanceRepository) UpdateBalance(ctx context.Context, userId string, amount float32) error {

	filter := bson.M{"userid": userId}
	update := bson.M{"$inc": bson.M{"amount": amount}}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *BalanceRepository) UpdateTotalCommission(ctx context.Context, userId string, amount float32) error {
	filter := bson.M{"userid": userId}
	update := bson.M{"$inc": bson.M{"total_commission": amount}}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *BalanceRepository) CreateBalance(ctx context.Context, userId string, amount float32) error {

	balance := types.Balance{
		UserId: userId,
		Amount: amount,
	}

	_, err := r.collection.InsertOne(ctx, balance)
	return err
}
