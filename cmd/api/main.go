package main

import (
	"context"
	"log"
	"os"

	"github.com/taketosaeki/donelog/internal/infrastructure/persistence/postgres"
)

func main() {
	ctx := context.Background()

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	pool, err := postgres.NewPool(ctx, connString)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}

	log.Println("database connection established (Supabase/Postgres). TODO: wire HTTP handlers and DI.")
}
