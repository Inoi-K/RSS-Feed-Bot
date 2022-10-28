package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/consts"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/env"
	"github.com/jackc/pgx/v5/pgconn"
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

	//query := fmt.Sprintf("SELECT EXISTS(SELECT 0 FROM schema1.users WHERE id = %v LIMIT 1);", userID)
	//var exists bool
	//err := db.pool.QueryRow(ctx, query).Scan(&exists)
	//if err != nil {
	//	return nil
	//}
	//if exists {
	//	return errors.New(fmt.Sprintf("user with id %v already exists", userID))
	//}

	query := fmt.Sprintf("INSERT INTO schema1.users VALUES (%v, 'ru');", userID)
	_, err := db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// TODO insert multiple rows in one query
func (db *Database) AddSource(ctx context.Context, userID int64, url string) error {
	query := fmt.Sprintf("INSERT INTO schema1.sources (url) VALUES ('%v') RETURNING id;", url)
	var sourceID int64
	err := db.pool.QueryRow(ctx, query).Scan(&sourceID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == consts.UniqueConstraintCode {
				query = fmt.Sprintf("SELECT id FROM schema1.sources WHERE url = '%v' LIMIT 1;", url)
				err := db.pool.QueryRow(ctx, query).Scan(&sourceID)
				if err != nil {
					return err
				}
			}
		} else {
			return err
		}
	}

	query = fmt.Sprintf("INSERT INTO schema1.userSource VALUES (%v, %v, true);", userID, sourceID)
	_, err = db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) RemoveSource(ctx context.Context, userID int64, url string) error {
	query := fmt.Sprintf("SELECT id FROM schema1.sources WHERE url = '%v' LIMIT 1;", url)
	var sourceID int64
	err := db.pool.QueryRow(ctx, query).Scan(&sourceID)
	if err != nil {
		return err
	}

	query = fmt.Sprintf("DELETE FROM schema1.userSource WHERE \"userId\" = %v AND \"sourceId\" = %v;", userID, sourceID)
	_, err = db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetUserSourcesTitleURL(ctx context.Context, userID int64) ([][]string, error) {
	query := fmt.Sprintf("SELECT title, url FROM schema1.sources WHERE id IN (SELECT \"sourceId\" FROM schema1.userSource WHERE \"userId\" = %v);", userID)
	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// TODO refactor []string{title, url} to a struct
	var sourceTitleURL [][]string
	for rows.Next() {
		var sourceTitle, sourceURL string
		err = rows.Scan(&sourceTitle, &sourceURL)
		if err != nil {
			return nil, err
		}
		sourceTitleURL = append(sourceTitleURL, []string{sourceTitle, sourceURL})
	}

	return sourceTitleURL, nil
}
