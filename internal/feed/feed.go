package feed

import (
	"context"
	"fmt"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/consts"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

var (
	cancel   func()
	canceled chan struct{}

	lastPostID int64
)

func Begin(ctx context.Context, bot *tgbotapi.BotAPI) {
	ticker := time.NewTicker(consts.FeedUpdateIntervalSeconds * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				ProcessNewPosts(ctx, bot)
			case <-canceled:
				ticker.Stop()
				return
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()

	// local cancel function
	canceled = make(chan struct{})
	cancel = func() {
		select {
		case <-canceled:
		default:
			close(canceled)
		}
	}
}

func End() {
	if cancel != nil {
		cancel()
	}
}

func tick(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "tick")
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("couldn't send tick: %v", err)
	}
}

// TODO implement new posts on chans & goroutines usage
func ProcessNewPosts(ctx context.Context, bot *tgbotapi.BotAPI) {
	db := database.GetDB()

	posts, err := db.GetNewPosts(ctx, lastPostID)
	if err != nil {
		log.Printf("couldn't get new posts: %v", err)
	}

	if len(posts) > 0 {
		lastPostID = posts[len(posts)-1].ID
	}

	for _, post := range posts {
		text := fmt.Sprintf("%v\n\n%v", post.Title, post.URL)
		msg := tgbotapi.NewMessage(post.ChatID, text)
		_, err := bot.Send(msg)
		if err != nil {
			log.Printf("couldn't send post in %v chat: %v", post.ChatID, err)
		}
	}

}
