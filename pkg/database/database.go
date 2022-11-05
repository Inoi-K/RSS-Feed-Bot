package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/consts"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/flags"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *Database

// Database provides ways for interaction with database (currently pool (proxy) only)
type Database struct {
	pool *pgxpool.Pool
}

// ConnectDB creates and returns connection to database
func ConnectDB(ctx context.Context) (*Database, error) {
	dbpool, err := pgxpool.New(ctx, *flags.DatabaseUrl)
	if err != nil {
		return nil, err
	}
	//defer dbpool.Close()

	db = &Database{
		pool: dbpool,
	}

	return db, nil
}

// GetDB returns database
func GetDB() *Database {
	return db
}

// AddUser creates a record of user in database
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

// AddSource adds one source in database and associates it with the user
func (db *Database) AddSource(ctx context.Context, userID int64, url string) error {
	query := fmt.Sprintf("INSERT INTO schema1.sources (url) VALUES ('%v') RETURNING id;", url)
	var sourceID int64
	err := db.pool.QueryRow(ctx, query).Scan(&sourceID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == consts.DuplicationCode {
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

// RemoveSource removes source user-source connection
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

// GetUserSourcesTitleURL gets title and url of the all sources associated with the user
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

// AlterSourceIsActive activates the source associated it with the user
func (db *Database) AlterSourceIsActive(ctx context.Context, chatID int64, url string, isActive bool) error {
	query := fmt.Sprintf("SELECT id FROM schema1.\"source\" WHERE url = '%v' LIMIT 1;", url)
	var sourceID int64
	err := db.pool.QueryRow(ctx, query).Scan(&sourceID)
	if err != nil {
		return err
	}

	query = fmt.Sprintf("UPDATE schema1.\"chatSource\" SET \"isActive\" = %v WHERE \"chatID\" = %v AND \"sourceID\" = %v;", isActive, chatID, sourceID)
	_, err = db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
