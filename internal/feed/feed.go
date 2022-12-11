package feed

import (
	"context"
	"fmt"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/consts"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/database"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/model"
	"github.com/Inoi-K/RSS-Feed-Bot/pkg/rss"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

var (
	cancel   func()
	canceled chan struct{}

	lastPostID     int64
	lastUpdateTime time.Time
)

// Begin initializes the ticker
func Begin(ctx context.Context, bot *tgbotapi.BotAPI) {
	ticker := time.NewTicker(consts.FeedUpdateIntervalSeconds * time.Second)
	lastUpdateTime = time.Now()

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

// End stops the ticker
func End() {
	if cancel != nil {
		cancel()
	}
}

// tick replies with "tick"
func tick(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "tick")
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("couldn't send tick: %v", err)
	}
}

// TODO implement new posts on chans & goroutines usage

// ProcessNewPosts handles getting and replying new posts
func ProcessNewPosts(ctx context.Context, bot *tgbotapi.BotAPI) {
	db := database.GetDB()
	beginTime := time.Now()

	urls, err := db.GetSourceURLs(ctx)
	if err != nil {
		log.Printf("couldn't get source urls: %v", err)
	}
	URLChat, err := db.GetSourceURLChat(ctx)
	if err != nil {
		log.Printf("couldn't get source urls of chats: %v", err)
	}

	// TODO split urls in several goroutines for parallelism
	for _, url := range urls {
		source, err := rss.Parse(url)
		if err != nil {
			log.Printf("couldn't parse %v link: %v", url, err)
		}

		posts := []model.Post{}
		for _, item := range source.Items {
			if item.PublishedParsed.Before(lastUpdateTime) {
				break
			}

			for _, chatID := range URLChat[url] {
				post := model.Post{
					Title:  item.Title,
					URL:    item.Link,
					ChatID: chatID,
				}
				posts = append(posts, post)
			}
		}
		// TODO delegate messaging to a worker in a goroutine
		for i := len(posts) - 1; i >= 0; i-- {
			post := posts[i]
			text := fmt.Sprintf("%v\n\n%v", post.Title, post.URL)
			msg := tgbotapi.NewMessage(post.ChatID, text)
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("couldn't send post in %v chat: %v", post.ChatID, err)
			}
		}
	}

	lastUpdateTime = beginTime
}
