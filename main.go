package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mmcdole/gofeed"
)

type Update struct {
	Title   string
	Link    string
	PubDate time.Time
}

func handleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message.Text == "/start" {
		sendMessage(bot, "Selamat datang! Silakan gunakan perintah lain untuk informasi lebih lanjut.")
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI("TOKEN")
	if err != nil {
		log.Fatal(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		handleMessage(bot, update)
	}

	fp := gofeed.NewParser()

	feedURL := "https://samehadaku.email/feed/"

	pollingInterval := time.Minute * 5

	var lastPubDate time.Time

	for {
		feed, err := fp.ParseURL(feedURL)
		if err != nil {
			log.Println("Error parsing feed:", err)
			continue
		}

		for _, item := range feed.Items {
			pubDate, err := time.Parse(time.RFC1123Z, item.Published)
			if err != nil {
				log.Println("Error parsing pubDate:", err)
				continue
			}

			if pubDate.After(lastPubDate) {
				sendMessage(bot, fmt.Sprintf("Judul Anime: %s\nLink Streaming: %s\nPosted by: <b>%s</b>", item.Title, item.Link, item.Author.Name))

				lastPubDate = pubDate
			}
		}

		time.Sleep(pollingInterval)
	}
}

func sendMessage(bot *tgbotapi.BotAPI, message string) {
	chatID := int64(CHAT_ID)
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "HTML"

	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Error sending message:", err)
	}
}
