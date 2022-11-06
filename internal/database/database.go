package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/consts"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/flags"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/structs"
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
	if db != nil {
		return db, nil
	}

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

// AddChat creates a record of chat in database
func (db *Database) AddChat(ctx context.Context, chatID int64, lang string) error {
	if len(lang) > 2 {
		return consts.LongLanguageError
	}

	//query := fmt.Sprintf("SELECT EXISTS(SELECT 0 FROM schema1.\"chat\" WHERE id = %v LIMIT 1);", chatID)
	//var exists bool
	//err := db.pool.QueryRow(ctx, query).Scan(&exists)
	//if err != nil {
	//	return nil
	//}
	//if exists {
	//	return errors.New(fmt.Sprintf("chat with id %v already exists", chatID))
	//}

	query := fmt.Sprintf("INSERT INTO schema1.\"chat\" VALUES (%v, '%v');", chatID, lang)
	_, err := db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// TODO insert multiple rows in one query

// AddSource adds one source in database and associates it with the chat
func (db *Database) AddSource(ctx context.Context, chatID int64, url string) error {
	query := fmt.Sprintf("INSERT INTO schema1.\"source\" (url) VALUES ('%v') RETURNING id;", url)
	var sourceID int64
	err := db.pool.QueryRow(ctx, query).Scan(&sourceID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == consts.DuplicationCode {
				query = fmt.Sprintf("SELECT id FROM schema1.\"source\" WHERE url = '%v' LIMIT 1;", url)
				err := db.pool.QueryRow(ctx, query).Scan(&sourceID)
				if err != nil {
					return err
				}
			}
		} else {
			return err
		}
	}

	query = fmt.Sprintf("INSERT INTO schema1.\"chatSource\" VALUES (%v, %v, true);", chatID, sourceID)
	_, err = db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// RemoveSource removes source chat-source connection
func (db *Database) RemoveSource(ctx context.Context, chatID int64, url string) error {
	query := fmt.Sprintf("SELECT id FROM schema1.\"source\" WHERE url = '%v' LIMIT 1;", url)
	var sourceID int64
	err := db.pool.QueryRow(ctx, query).Scan(&sourceID)
	if err != nil {
		return err
	}

	query = fmt.Sprintf("DELETE FROM schema1.\"chatSource\" WHERE \"chatID\" = %v AND \"sourceID\" = %v;", chatID, sourceID)
	_, err = db.pool.Query(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// GetChatSourcesTitleURL gets title and url of the all sources associated with the chat
func (db *Database) GetChatSourcesTitleURL(ctx context.Context, chatID int64) ([][]string, error) {
	query := fmt.Sprintf("SELECT title, url FROM schema1.\"source\" WHERE id IN (SELECT \"sourceID\" FROM schema1.\"chatSource\" WHERE \"chatID\" = %v);", chatID)
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

// GetNewPosts returns slice of the posts with id greater than the most recent post id
func (db *Database) GetNewPosts(ctx context.Context, lastPostID int64) ([]structs.Post, error) {
	query := fmt.Sprintf("SELECT * FROM schema1.post WHERE id > %v;", lastPostID)
	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []structs.Post
	for rows.Next() {
		var post structs.Post
		err = rows.Scan(&post.ID, &post.SourceID, &post.Title, &post.URL, &post.ChatID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// TODO generalize with GetUserSourcesTitleURL

// GetChatSourcesTitleURLByIsActive gets title and url of the all active/deactive sources associated with the user
func (db *Database) GetChatSourcesTitleURLByIsActive(ctx context.Context, userID int64, isActive bool) ([][]string, error) {
	query := fmt.Sprintf("SELECT title, url FROM schema1.source WHERE id IN (SELECT \"sourceID\" FROM schema1.\"chatSource\" WHERE \"chatID\" = %v AND \"isActive\" = %v);", userID, isActive)
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
