package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("All"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Research swtch"),
		tgbotapi.NewKeyboardButton("Habr"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Russia Today"),
		tgbotapi.NewKeyboardButton("Lenta ru"),
		tgbotapi.NewKeyboardButton("New York Times"),
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

			fmt.Println(update.Message.Command())
			fmt.Println(update.Message.Text)

			if update.Message.IsCommand() {
				fmt.Printf("#####################:%d:\n", 1)
				if update.Message.Command() == "menu" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose your selection")
					msg.ReplyMarkup = numericKeyboard
					bot.Send(msg)
				}
			} else {
				switch update.Message.Text {
				case "All":
					// print all news
				case "Research swtch":
					// only research swtch
				case "Habr":
					// only habr news
				case "Russia Today":
					// only RT
				case "Lenta ru":
					// only Lenta ru news
				case "New York Times":
					// only NYT
				default:
					fmt.Printf("#####################:%d:\n", 2)
					ansText := "Wow, I'm sorry," +
						" but I was created only for sending news" +
						" not for conversation:("
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, ansText)
					bot.Send(msg)
				}
			}

			// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			// msg.ReplyToMessageID = update.Message.MessageID

		}
	}
}
