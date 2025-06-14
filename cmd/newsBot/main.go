package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/SANEKNAYMCHIK/newsBot/internal/app/parser"
	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	_ "net/http/pprof"

	"github.com/SANEKNAYMCHIK/newsBot/proto/news"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
			tgbotapi.NewKeyboardButton("Research-swtch"),
			tgbotapi.NewKeyboardButton("Habr"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Russia-Today"),
			tgbotapi.NewKeyboardButton("Lenta-ru"),
			tgbotapi.NewKeyboardButton("New-York-Times"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Close keyboard"),
			tgbotapi.NewKeyboardButton("Update news"),
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

type NewsStorage struct {
	mu   sync.RWMutex
	news map[string][]services.NewsItem
}

func getNewsFromParser(sources []string) (*map[string]*news.NewsItemList, error) {
	// conn, err := grpc.NewClient("newsParser:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := news.NewNewsParserClient(conn)
	resp, err := client.GetNews(context.Background(), &news.NewsRequest{Sources: sources})
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func convertProtoToServices(data *map[string]*news.NewsItemList) map[string][]services.NewsItem {
	result := make(map[string][]services.NewsItem)
	for name, itemList := range *data {
		for _, item := range itemList.Items {
			time, err := time.Parse("2006-01-02 15:04:05", item.Date)
			if err != nil {
				log.Printf("Error parsing date %s for item %s: %v", item.Date, item.Title, err)
			}
			result[name] = append(result[name], services.NewsItem{
				Title:       item.Title,
				Link:        item.Link,
				Date:        time,
				Description: item.Description,
				Website:     item.Website,
			})
		}
	}
	return result
}

func main() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()
	sources := []string{
		"https://habr.com/ru/rss/articles/",
		"https://russian.rt.com/rss",
		"https://lenta.ru/rss",
		"https://nytimes.com/services/xml/rss/nyt/World.xml",
		"https://research.swtch.com/feed.atom",
	}

	// storage := &NewsStorage{
	// 	news: parser.ParseAllNews(sources),
	// }
	newsData, err := getNewsFromParser(sources)
	if err != nil {
		log.Fatal("Failed to fetch news from parser")

	}
	storage := &NewsStorage{
		news: convertProtoToServices(newsData),
	}

	godotenv.Load()
	token := os.Getenv("TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Printf("Some errors with tgbot: %s", err)
		panic(err)
	}

	setCommands(bot)

	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		go func(update tgbotapi.Update) {
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
					case "Research-swtch", "Habr", "Russia-Today", "Lenta-ru", "New-York-Times":
						storage.mu.RLock()
						newsVals := storage.news[update.Message.Text]
						storage.mu.RUnlock()
						for i := 0; i < len(newsVals); i++ {
							ansText := ""
							// fmt.Println(t.Format("2006-01-02 15:04:05"))
							// newsVals[update.Message.Text][i].Date.Format("2006-01-02 15:04:05")
							ansText += newsVals[i].Date.String() + "\n"
							ansText += newsVals[i].Title + "\n"
							ansText += newsVals[i].Description + "\n"
							ansText += newsVals[i].Link
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, ansText)
							bot.Send(msg)
							time.Sleep(1 * time.Second)
						}
					case "Close keyboard":
						hideKeyboard(bot, update.Message.Chat.ID, update.Message.MessageID)
					case "Update news":
						go func(chatID int64) {
							msg := tgbotapi.NewMessage(chatID, "Обновляем новости")
							bot.Send(msg)
							latestNews := parser.ParseAllNews(sources)
							storage.mu.Lock()
							storage.news = latestNews
							storage.mu.Unlock()
							msg = tgbotapi.NewMessage(chatID, "Новости обновлены")
							bot.Send(msg)
						}(update.Message.Chat.ID)
					default:
						ansText := "Wow, I'm sorry," +
							" but I was created only for sending news" +
							" not for conversation:("
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, ansText)
						bot.Send(msg)
					}
				}
			}
		}(update)
	}
}
