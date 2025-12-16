package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/SANEKNAYMCHIK/newsBot/internal/bot"
	"github.com/SANEKNAYMCHIK/newsBot/internal/config"
	"github.com/SANEKNAYMCHIK/newsBot/internal/database"
	"github.com/SANEKNAYMCHIK/newsBot/internal/repositories"
	"github.com/SANEKNAYMCHIK/newsBot/internal/services"
	"github.com/SANEKNAYMCHIK/newsBot/pkg/auth"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	db, err := database.NewPostgres(ctx, cfg.DBUrl)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	log.Println("Подключение к БД установлено")

	userRepo := repositories.NewUserRepository(db.Pool)
	sourceRepo := repositories.NewSourceRepository(db.Pool)
	newsRepo := repositories.NewNewsRepository(db.Pool)
	subscriptionRepo := repositories.NewSubscriptionRepository(db.Pool)
	categoryRepo := repositories.NewCategoryRepository(db.Pool)

	jwtManager := auth.NewJWTManager(cfg.JWTSecret)
	authService := services.NewAuthService(userRepo, jwtManager)
	adminService := services.NewAdminService(userRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	sourceService := services.NewSourceService(sourceRepo)
	newsService := services.NewNewsService(newsRepo, sourceRepo, subscriptionRepo)

	rssParser := services.NewRssParser(10)
	rssService := services.NewRssService(sourceRepo, newsRepo, rssParser)
	refreshService := services.NewRefreshService(
		rssService,
		subscriptionRepo,
		5,
		100,
		3*time.Minute,
	)
	go refreshService.Start(context.Background())

	botService := bot.NewBotService(
		authService,
		sourceRepo,
		userRepo,
		newsRepo,
		subscriptionRepo,
		categoryRepo,
		newsService,
		adminService,
		categoryService,
		sourceService,
		refreshService,
	)

	telegramBot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatalf("Ошибка создания бота: %v", err)
	}

	telegramBot.Debug = false
	log.Printf("Авторизован как бот: %s", telegramBot.Self.UserName)

	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Начать работу с ботом"},
		{Command: "help", Description: "Показать помощь"},
		{Command: "subscribe", Description: "Управление подписками"},
		{Command: "news", Description: "Последние новости"},
		{Command: "source_news", Description: "Новости конкретного источника"},
		{Command: "sources", Description: "Доступные источники"},
		{Command: "add_source", Description: "Добавить новый источник"},
		{Command: "categories", Description: "Показать категории"},
		{Command: "update", Description: "Обновить новости вручную"},
		{Command: "update_status", Description: "Статус обновления новостей"},
	}

	adminCommands := []tgbotapi.BotCommand{
		{Command: "admin", Description: "Админ-панель"},
		{Command: "admin_users", Description: "Список пользователей"},
		{Command: "admin_stats", Description: "Статистика системы"},
		{Command: "admin_make_admin", Description: "Назначить админа"},
		{Command: "admin_remove_admin", Description: "Снять админа"},
		{Command: "admin_add_category", Description: "Добавить категорию"},
		{Command: "admin_update_source", Description: "Изменить источник"},
	}

	allCommands := append(commands, adminCommands...)

	if _, err := telegramBot.Request(tgbotapi.NewSetMyCommands(allCommands...)); err != nil {
		log.Printf("Ошибка настройки команд: %v", err)
	}

	handler := bot.NewHandler(telegramBot, botService)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := telegramBot.GetUpdatesChan(u)

	log.Println("Бот запущен. Нажмите Ctrl+C для остановки.")

	sem := make(chan struct{}, 50)
	var wg sync.WaitGroup

	for {
		select {
		case <-quit:
			log.Println("Остановка бота...")
			telegramBot.StopReceivingUpdates()
			wg.Wait()
			log.Println("Бот остановлен")
			return
		case update := <-updates:
			select {
			case sem <- struct{}{}:
				wg.Add(1)
				go func(update tgbotapi.Update) {
					defer wg.Done()
					defer func() {
						<-sem
					}()
					handler.HandleUpdate(update)
				}(update)
			default:
				log.Printf("Все занято, пропускаем обновление: %v", update)
			}
		}
	}
}
