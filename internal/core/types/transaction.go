package types

import "time"

type TransactionType string

const (
	DEPOSIT    TransactionType = "DEPOSIT"
	PURCHASE   TransactionType = "PURCHASE"
	COMMISSION TransactionType = "COMMISSION"
)

type Transaction struct {
	Id                 string          `bson:"id"`
	Type               TransactionType `bson:"type"`
	Amount             float32         `bson:"amount"`
	Customer           User            `bson:"customer"`
	Restaurant         User            `bson:"restaurant"`
	CreatedAt          time.Time       `bson:"created_at"`
	RelatedTransaction string          `bson:"related_transaction"`
}
