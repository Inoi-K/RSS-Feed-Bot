package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/env"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbpool *pgxpool.Pool
)

func main() {
	ctx := context.Background()

	var err error
	dbpool, err = pgxpool.New(ctx, *env.DatabaseUrl)
	if err != nil {
		fmt.Printf("unable to connect to database: %v\n", err)
		return
	}
	defer dbpool.Close()

	err = add(6, "cc")
	if err != nil {
		fmt.Println(err)
	}
}

func add(id int, lang string) error {
	if len(lang) > 2 {
		return errors.New("language can't be longer than 2 symbols")
	}

	query := fmt.Sprintf("INSERT INTO schema1.users (id, lang) VALUES (%v, '%v');", id, lang)
	_, err := dbpool.Query(context.Background(), query)
	if err != nil {
		return err
	}

	return nil
}
