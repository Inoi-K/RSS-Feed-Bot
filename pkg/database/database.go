package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/env"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *Database

type Database struct {
	pool *pgxpool.Pool
}

func ConnectDB(ctx context.Context) (*Database, error) {
	dbpool, err := pgxpool.New(ctx, *env.DatabaseUrl)
	if err != nil {
		return nil, err
	}
	//defer dbpool.Close()

	db = &Database{
		pool: dbpool,
	}

	return db, nil
}

func GetDB() *Database {
	return db
}

func (db *Database) AddUser(ctx context.Context, userID int64, lang string) error {
	if len(lang) > 2 {
		return errors.New("language cannot be longer than 2 symbols")
	}

	query := fmt.Sprintf("SELECT EXISTS(SELECT 0 FROM schema1.users WHERE id = %v LIMIT 1);", userID)
	var exists bool
	err := db.pool.QueryRow(ctx, query).Scan(&exists)
	if err != nil {
		return nil
	}
	if exists {
		return errors.New(fmt.Sprintf("user with id %v already exists", userID))
	}

	query = fmt.Sprintf("INSERT INTO schema1.users (id, lang) VALUES (%v, '%v');", userID, lang)
	_, err = db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
