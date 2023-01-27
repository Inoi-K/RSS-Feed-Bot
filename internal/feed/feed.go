package feed

import (
	"context"
	"fmt"
	"github.com/Inoi-K/RSS-Feed-Bot/configs/flags"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/builder"
	db "github.com/Inoi-K/RSS-Feed-Bot/internal/database"
	"github.com/Inoi-K/RSS-Feed-Bot/internal/model"
	"github.com/Inoi-K/RSS-Feed-Bot/pkg/parser"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

var (
	cancel   func()
	canceled chan struct{}

	lastUpdateTime time.Time
)

// Begin initializes the ticker
func Begin(ctx context.Context, bot *tgbotapi.BotAPI) {
	ticker := time.NewTicker(*flags.FeedUpdateInterval)
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
	msg := builder.NewMessage(chatID, "tick", nil)
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("couldn't send tick: %v", err)
	}
}

// ProcessNewPosts handles getting and replying new posts
func ProcessNewPosts(ctx context.Context, bot *tgbotapi.BotAPI) {
	beginTime := time.Now()

	URLChat, err := db.GetSourceURLChat(ctx)
	if err != nil {
		log.Printf("couldn't get source urls of chats: %v", err)
	}

	// chan for posts
	posts := make(chan model.Post, 50)
	done := make(chan struct{})

	// parse url and send its posts to the chan
	parse := func(url string, chats []int64) {
		source, err := parser.Parse(url)
		if err != nil {
			log.Printf("couldn't parse %v link: %v", url, err)
		}

		for _, item := range source.Items {
			if item.PublishedParsed == nil || item.PublishedParsed.Before(lastUpdateTime) {
				continue
			}

			for _, chatID := range chats {
				post := model.Post{
					Title:  item.Title,
					URL:    item.Link,
					ChatID: chatID,
				}
				posts <- post
			}
		}

		done <- struct{}{}
	}

	// send an incoming from the chan post to telegram
	sendMessage := func() {
		for post := range posts {
			text := fmt.Sprintf("%v\n\n%v", post.Title, post.URL)
			msg := builder.NewMessage(post.ChatID, text, nil)
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("couldn't send post in %v chat: %v", post.ChatID, err)
			}
		}
	}

	// launch 5 'send' workers
	go func() {
		for i := 0; i < 5; i++ {
			go sendMessage()
		}
	}()

	// process all urls
	for url, chats := range URLChat {
		go parse(url, chats)
	}

	// waiting until all posts are parsed and stacked for sending
	for i := 0; i < len(URLChat); i++ {
		<-done
	}
	close(posts)

	lastUpdateTime = beginTime
}
