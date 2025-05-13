package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/SANEKNAYMCHIK/newsBot/internal/app/parser"
	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func setCommands(bot *tgbotapi.BotAPI) {
	commands := []tgbotapi.BotCommand{
		{
			Command:     "start",
			Description: "Начать работу с ботом",
		},
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
	menu.ResizeKeyboard = true

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

	newsVals := make(map[string][]services.NewsItem)

	for item := range ch {
		newsVals[item.Website] = append(newsVals[item.Website], item)
	}

	godotenv.Load()
	token := os.Getenv("TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Printf("Some errors with tgbot: %s", err)
		panic(err)
	}

	// bot.Debug = true
	setCommands(bot)

	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {

			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "news":
					showMainMenu(bot, update.Message.Chat.ID)
				case "start":
					ansText := "Hi! I can help to you with the latest news\n" +
						"Click on a news button."
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, ansText)
					bot.Send(msg)
				}
			} else {
				switch update.Message.Text {
				case "All":
					for j := range newsVals {
						for i := 0; i < len(newsVals[j]); i++ {
							ansText := ""
							// fmt.Println(t.Format("2006-01-02 15:04:05"))
							// newsVals[update.Message.Text][i].Date.Format("2006-01-02 15:04:05")
							ansText += newsVals[j][i].Date.String() + "\n"
							ansText += newsVals[j][i].Title + "\n"
							ansText += newsVals[j][i].Description + "\n"
							ansText += newsVals[j][i].Link
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, ansText)
							bot.Send(msg)
							time.Sleep(1 * time.Second)
						}
					}
				case "Research-swtch":
					for i := 0; i < len(newsVals[update.Message.Text]); i++ {
						ansText := ""
						// fmt.Println(t.Format("2006-01-02 15:04:05"))
						// newsVals[update.Message.Text][i].Date.Format("2006-01-02 15:04:05")
						ansText += newsVals[update.Message.Text][i].Date.String() + "\n"
						ansText += newsVals[update.Message.Text][i].Title + "\n"
						ansText += newsVals[update.Message.Text][i].Description + "\n"
						ansText += newsVals[update.Message.Text][i].Link
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, ansText)
						bot.Send(msg)
						time.Sleep(1 * time.Second)
					}
				case "Habr":
					for i := 0; i < len(newsVals[update.Message.Text]); i++ {
						ansText := ""
						// fmt.Println(t.Format("2006-01-02 15:04:05"))
						// newsVals[update.Message.Text][i].Date.Format("2006-01-02 15:04:05")
						ansText += newsVals[update.Message.Text][i].Date.String() + "\n"
						ansText += newsVals[update.Message.Text][i].Title + "\n"
						ansText += newsVals[update.Message.Text][i].Description + "\n"
						ansText += newsVals[update.Message.Text][i].Link
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, ansText)
						bot.Send(msg)
						time.Sleep(1 * time.Second)
					}
				case "Russia-Today":
					for i := 0; i < len(newsVals[update.Message.Text]); i++ {
						ansText := ""
						// fmt.Println(t.Format("2006-01-02 15:04:05"))
						// newsVals[update.Message.Text][i].Date.Format("2006-01-02 15:04:05")
						ansText += newsVals[update.Message.Text][i].Date.String() + "\n"
						ansText += newsVals[update.Message.Text][i].Title + "\n"
						ansText += newsVals[update.Message.Text][i].Description + "\n"
						ansText += newsVals[update.Message.Text][i].Link
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, ansText)
						bot.Send(msg)
						time.Sleep(1 * time.Second)
					}
				case "Lenta-ru":
					for i := 0; i < len(newsVals[update.Message.Text]); i++ {
						ansText := ""
						// fmt.Println(t.Format("2006-01-02 15:04:05"))
						// newsVals[update.Message.Text][i].Date.Format("2006-01-02 15:04:05")
						ansText += newsVals[update.Message.Text][i].Date.String() + "\n"
						ansText += newsVals[update.Message.Text][i].Title + "\n"
						ansText += newsVals[update.Message.Text][i].Description + "\n"
						ansText += newsVals[update.Message.Text][i].Link
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, ansText)
						bot.Send(msg)
						time.Sleep(1 * time.Second)
					}
				case "New-York-Times":
					for i := 0; i < len(newsVals[update.Message.Text]); i++ {
						ansText := ""
						// fmt.Println(t.Format("2006-01-02 15:04:05"))
						// newsVals[update.Message.Text][i].Date.Format("2006-01-02 15:04:05")
						ansText += newsVals[update.Message.Text][i].Date.String() + "\n"
						ansText += newsVals[update.Message.Text][i].Title + "\n"
						ansText += newsVals[update.Message.Text][i].Description + "\n"
						ansText += newsVals[update.Message.Text][i].Link
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, ansText)
						bot.Send(msg)
						time.Sleep(1 * time.Second)
					}
				case "Close keyboard":
					hideKeyboard(bot, update.Message.Chat.ID, update.Message.MessageID)
				default:
					ansText := "Wow, I'm sorry," +
						" but I was created only for sending news" +
						" not for conversation:("
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, ansText)
					bot.Send(msg)
				}
			}
		}
	}
}
