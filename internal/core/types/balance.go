package types

type Balance struct {
	UserId          string  `bson:"userid"`
	Amount          float32 `bson:"amount"`
	TotalCommission float32 `bson:"total_commission"`
}