package types

type UserType string

const (
	CUSTOMER   UserType = "CUSTOMER"
	RESTAURANT UserType = "RESTAURANT"
)

type User struct {
	Id   string   `bson:"id"`
	Type UserType `bson:"type"`
}