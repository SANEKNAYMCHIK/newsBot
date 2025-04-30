package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("All", "1"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Research swtch", "2"),
		tgbotapi.NewInlineKeyboardButtonData("Habr", "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Russia Today", "4"),
		tgbotapi.NewInlineKeyboardButtonData("Lenta ru", "5"),
		tgbotapi.NewInlineKeyboardButtonData("New York Times", "6"),
	),
)

func main() {
	// fmt.Println(bot)

	// sources := []string{
	// 	"https://habr.com/ru/rss/articles/",
	// 	"https://russian.rt.com/rss",
	// 	"https://lenta.ru/rss",
	// 	"https://rss.nytimes.com/services/xml/rss/nyt/World.xml",
	// 	"https://research.swtch.com/feed.atom",
	// }
	// var wg sync.WaitGroup
	// // channel for keeping news
	// ch := make(chan services.NewsItem, 100)
	// for _, url := range sources {
	// 	wg.Add(1)
	// 	go parser.Parse(url, &wg, ch)
	// }
	// go func() {
	// 	wg.Wait()
	// 	close(ch)
	// }()
	// var news []services.NewsItem
	// for item := range ch {
	// 	news = append(news, item)
	// 	// fmt.Println(item.Link)
	// }

	godotenv.Load()
	token := os.Getenv("TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Printf("Some errors with tgbot: %s", err)
		panic(err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			fmt.Println(update.Message.Command())
			fmt.Println(update.Message.Text)

			if update.Message.Command() == "menu" {
				keyb := tgbotapi.InlineKeyboardMarkup{}
				var numkey []tgbotapi.InlineKeyboardButton
				numkey = append(numkey, tgbotapi.NewInlineKeyboardButtonData("2", "2"))
				keyb.InlineKeyboard = append(keyb.InlineKeyboard, numkey)
				numkey = numkey[:0]
				numkey = append(numkey, tgbotapi.NewInlineKeyboardButtonData("3", "3"))
				keyb.InlineKeyboard = append(keyb.InlineKeyboard, numkey)
				msg.ReplyMarkup = numericKeyboard
			}
			// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			// msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}
