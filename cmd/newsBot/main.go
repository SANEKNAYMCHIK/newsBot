package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/SANEKNAYMCHIK/newsBot/internal/app/parser"
	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func setCommands(bot *tgbotapi.BotAPI) {
	commands := []tgbotapi.BotCommand{
		{
			Command:     "start",
			Description: "Начать работу с ботом",
		},
		// {
		// 	Command:     "help",
		// 	Description: "Показать справку",
		// },
		{
			Command:     "news",
			Description: "Показать новостные порталы",
		},
	}
	config := tgbotapi.NewSetMyCommands(commands...)
	_, err := bot.Request(config)
	if err != nil {
		log.Panic(err)
	}
}

func showMainMenu(bot *tgbotapi.BotAPI, chatID int64) {
	// Создаем клавиатуру меню
	menu := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("All"),
			tgbotapi.NewKeyboardButton("Close keyboard"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Research-swtch"),
			tgbotapi.NewKeyboardButton("Habr"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Russia-Today"),
			tgbotapi.NewKeyboardButton("Lenta-ru"),
			tgbotapi.NewKeyboardButton("New-York-Times"),
		),
	)
	menu.ResizeKeyboard = true // делает кнопки компактнее

	msg := tgbotapi.NewMessage(chatID, "Выберите портал:")
	msg.ReplyMarkup = menu
	bot.Send(msg)
}

func hideKeyboard(bot *tgbotapi.BotAPI, chatID int64, msgID int) {
	bot.Send(tgbotapi.NewDeleteMessage(chatID, msgID))
	msg := tgbotapi.NewMessage(chatID, "Keyboard is closed")
	msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
	bot.Send(msg)
}

func main() {
	sources := []string{
		"https://habr.com/ru/rss/articles/",
		"https://russian.rt.com/rss",
		"https://lenta.ru/rss",
		"https://nytimes.com/services/xml/rss/nyt/World.xml",
		"https://research.swtch.com/feed.atom",
	}
	var wg sync.WaitGroup
	// channel for keeping news
	ch := make(chan services.NewsItem, 100)

	for _, url := range sources {
		wg.Add(1)
		go parser.Parse(url, &wg, ch)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	// newsVals := make(map[string][]services.NewsItem)

	// var news []services.NewsItem
	for item := range ch {
		// newsVals[item]
		// news = append(news, item)
		fmt.Println(item.Website)
	}

	// godotenv.Load()
	// token := os.Getenv("TOKEN")
	// bot, err := tgbotapi.NewBotAPI(token)

	// if err != nil {
	// 	log.Printf("Some errors with tgbot: %s", err)
	// 	panic(err)
	// }

	// // bot.Debug = true
	// setCommands(bot)

	// log.Printf("Authorized on account %s", bot.Self.UserName)

	// u := tgbotapi.NewUpdate(0)
	// u.Timeout = 60

	// updates := bot.GetUpdatesChan(u)

	// for update := range updates {
	// 	if update.Message != nil {
	// 		fmt.Println(update.Message.Command())
	// 		fmt.Println(update.Message.Text)
	// 		fmt.Println(update.Message.IsCommand())

	// 		if update.Message.IsCommand() {
	// 			fmt.Printf("#####################:%d:\n", 1)
	// 			switch update.Message.Command() {
	// 			case "news":
	// 				showMainMenu(bot, update.Message.Chat.ID)
	// 			case "start":
	// 				ansText := "Hi! I can help to you with the latest news\n" +
	// 					"Click on a news button."
	// 				msg := tgbotapi.NewMessage(update.Message.Chat.ID, ansText)
	// 				bot.Send(msg)
	// 			}
	// 		} else {
	// 			switch update.Message.Text {
	// 			case "All":
	// 				// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	// 				// bot.Send(msg)
	// 			case "Research-swtch":
	// 				// only research swtch
	// 			case "Habr":
	// 				// only habr news
	// 			case "Russia-Today":
	// 				// only RT
	// 			case "Lenta-ru":
	// 				// only Lenta ru news
	// 			case "New-York-Times":
	// 				// only NYT
	// 			case "Close keyboard":
	// 				hideKeyboard(bot, update.Message.Chat.ID, update.Message.MessageID)
	// 			default:
	// 				ansText := "Wow, I'm sorry," +
	// 					" but I was created only for sending news" +
	// 					" not for conversation:("
	// 				msg := tgbotapi.NewMessage(update.Message.Chat.ID, ansText)
	// 				bot.Send(msg)
	// 			}
	// 		}

	// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	// msg.ReplyToMessageID = update.Message.MessageID

	// }
	// }
}
