package database

import (
	"context"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/consts"
	"github.com/jackc/pgx/v5/pgconn"
	"testing"
)

func TestConnectDB(t *testing.T) {

}

func TestAddChat(t *testing.T) {
	tests := []struct {
		name   string
		chatID int64
		lang   string
		want   error
	}{
		{"lang >2 symbols", 1, "eng", consts.LongLanguageError},
		{"new chat", 1, "en", nil},
		{"existing chat", 1, "ru", &pgconn.PgError{}},
	}

	ctx := context.Background()
	ConnectDB(ctx)
	db := GetDB()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := db.AddChat(ctx, test.chatID, test.lang)
			if got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}

func TestAddSource(t *testing.T) {
	tests := []struct {
		name   string
		chatID int64
		url    string
		want   error
	}{
		{"new source", 1, "source1", nil},
		{"existing source", 1, "source1", nil},
		{"same source", 2, "source1", nil},
	}

	ctx := context.Background()
	ConnectDB(ctx)
	db := GetDB()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := db.AddSource(ctx, test.chatID, test.url)
			if got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}

func TestRemoveSource(t *testing.T) {
	tests := []struct {
		name   string
		chatID int64
		url    string
		want   error
	}{
		{"existing source", 1, "source1", nil},
		{"non-existing source", 1, "no-source", nil},
	}

	ctx := context.Background()
	ConnectDB(ctx)
	db := GetDB()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := db.AddSource(ctx, test.chatID, test.url)
			if got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}

func TestGetChatSourcesTitleURL(t *testing.T) {
	tests := []struct {
		name   string
		chatID int64
		want   error
	}{
		{"existing chatID", 1, nil},
		{"non-existing chatID", -1, nil},
	}

	ctx := context.Background()
	ConnectDB(ctx)
	db := GetDB()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, got := db.GetChatSourcesTitleURL(ctx, test.chatID)
			if got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}

func TestGetNewPosts(t *testing.T) {
	tests := []struct {
		name       string
		lastPostID int64
		want       error
	}{
		{"negative lastPostID", -1, nil},
		{"existing lastPostID", 1, nil},
		{"out of bounds lastPostID", 1_000_000_000, nil},
	}

	ctx := context.Background()
	ConnectDB(ctx)
	db := GetDB()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, got := db.GetNewPosts(ctx, test.lastPostID)
			if got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}
