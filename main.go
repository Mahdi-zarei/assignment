package main

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func main() {
	db, err := pgx.Connect(context.Background(), "postgres://postgres:dummypass@5.34.202.174:5433/giftshop")
	if err != nil {
		panic(err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	db.Close(context.Background())
}
