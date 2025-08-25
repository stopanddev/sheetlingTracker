package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var Conn *pgxpool.Pool

func Init() {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		log.Fatal("DATABASE_URL not set")
	}

	var err error
	Conn, err = pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
}
