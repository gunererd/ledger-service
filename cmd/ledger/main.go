package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ledger-service/internal/core/services/ledger"
	"ledger-service/internal/infrastructure/config"
	"ledger-service/internal/infrastructure/db"
	"ledger-service/internal/infrastructure/queue"
	"ledger-service/internal/infrastructure/repository/mongo"
	"ledger-service/internal/infrastructure/web"
)

func main() {
	cfg := config.LoadFromEnv()

	client, err := db.NewMongoClient(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	transactionRepo := mongo.NewTransactionRepository(client, cfg.DatabaseName, cfg.TransactionCollection)
	balanceRepo := mongo.NewBalanceRepository(client, cfg.DatabaseName, cfg.BalanceCollection)
	taskQueue := queue.NewInMemoryQueue()

	ledgerService := ledger.NewService(transactionRepo, balanceRepo, taskQueue)

	server := web.NewServer(ledgerService)

	addr := ":" + cfg.ServerPort
	fmt.Printf("Server starting on %s\n", addr)
	fmt.Printf("API documentation available at http://localhost%s/docs\n", addr)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: server.Handler(),
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	gracefulShutdown(httpServer, ledgerService)
}

func gracefulShutdown(httpServer *http.Server, ledgerService *ledger.Service) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	ledgerService.Shutdown()

	fmt.Println("Server stopped")
}
