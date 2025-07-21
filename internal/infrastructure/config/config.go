package config

import "os"

type Config struct {
	MongoURI              string
	DatabaseName          string
	TransactionCollection string
	BalanceCollection     string
	ServerPort            string
	QueueBufferSize       int
}

func LoadFromEnv() *Config {
	return &Config{
		MongoURI:              getEnv("MONGO_URI", "mongodb://admin:password@localhost:27017/ledger?authSource=admin"),
		DatabaseName:          getEnv("DATABASE_NAME", "ledger"),
		TransactionCollection: getEnv("TRANSACTION_COLLECTION", "transactions"),
		BalanceCollection:     getEnv("BALANCE_COLLECTION", "balances"),
		ServerPort:            getEnv("SERVER_PORT", "8081"),
		QueueBufferSize:       100,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
