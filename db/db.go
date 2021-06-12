package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var connectionString = "postgres://amqjqcvlywltns:49b3cea25bb797e1e7cc4884702d6e7e08558f08e9ffc49cc086a8aecd6c219c@ec2-54-73-68-39.eu-west-1.compute.amazonaws.com:5432/d4143he11vfch2"

// SetDB to configure db connection
func SetDB() (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return dbpool, nil
}
