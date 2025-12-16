package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot     *tgbotapi.BotAPI
	service *BotService
}

func NewHandler(bot *tgbotapi.BotAPI, service *BotService) *Handler {
	return &Handler{
		bot:     bot,
		service: service,
	}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update) {
	ctx := context.Background()

	if update.Message != nil {
		h.handleMessage(ctx, update.Message)
	} else if update.CallbackQuery != nil {
		h.handleCallbackQuery(ctx, update.CallbackQuery)
	}
}

func (h *Handler) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	h.bot.Send(msg)
}

func (h *Handler) handleMessage(ctx context.Context, message *tgbotapi.Message) {
	if message.From == nil {
		return
	}
	user, err := h.service.authService.RegisterOrUpdateTelegramUser(
		ctx,
		message.From.ID,
		message.From.UserName,
		message.From.FirstName+" "+message.From.LastName,
	)
	if err != nil {
		log.Printf("Error user's register: %v", err)
		h.sendMessage(message.Chat.ID, "Ошибка регистрации. Попробуйте еще раз.")
		return
	}
	if message.IsCommand() {
		h.handleCommand(ctx, message, user)
		return
	}

	h.handleText(ctx, message, user)
}

func (h *Handler) handleCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	switch message.Command() {
	case "start":
		h.handleStart(ctx, message, user)
	case "help":
		h.handleHelp(ctx, message)
	case "subscribe":
		h.handleSubscribeCommand(ctx, message, user)
	case "news":
		h.handleNewsCommand(ctx, message, user)
	case "sources":
		h.handleSourcesCommand(ctx, message, user)
	case "source_news":
		h.handleSourceNewsCommand(ctx, message, user)
	case "add_source":
		h.handleAddSourceCommand(ctx, message, user)
	case "categories":
		h.handleCategoriesCommand(ctx, message)
	case "admin":
		h.handleAdminCommand(ctx, message, user)
	case "admin_users":
		h.handleAdminUsersCommand(ctx, message, user)
	case "admin_stats":
		h.handleAdminStatsCommand(ctx, message, user)
	case "admin_make_admin":
		h.handleAdminMakeAdminCommand(ctx, message, user)
	case "admin_remove_admin":
		h.handleAdminRemoveAdminCommand(ctx, message, user)
	case "admin_add_category":
		h.handleAdminAddCategoryCommand(ctx, message, user)
	case "admin_update_source":
		h.handleAdminUpdateSourceCommand(ctx, message, user)
	case "update":
		h.handleUpdateCommand(ctx, message, user)
	case "update_status":
		h.handleUpdateStatusCommand(ctx, message, user)
	default:
		h.sendMessage(message.Chat.ID, "Неизвестная команда. Используйте /help для списка команд.")
	}
}

func (h *Handler) handleStart(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	welcomeText := fmt.Sprintf(
		"Привет, %s! Я — новостной бот.\n\n"+
			"Я могу:\n"+
			"• Подписать вас на желаемый новостной источник\n"+
			"• Присылать свежие новости\n"+
			"• Показывать новости по категориям\n\n"+
			"Используйте команды или кнопки ниже:",
		*user.TgFirstName,
	)

	msg := tgbotapi.NewMessage(message.Chat.ID, welcomeText)
	msg.ReplyMarkup = MainMenuKeyboard()
	h.bot.Send(msg)
}

func (h *Handler) handleHelp(ctx context.Context, message *tgbotapi.Message) {
	isAdmin := false
	user, err := h.service.authService.RegisterOrUpdateTelegramUser(
		ctx,
		message.From.ID,
		message.From.UserName,
		message.From.FirstName+" "+message.From.LastName,
	)
	if err == nil {
		if adminCheck, err := h.service.IsAdmin(ctx, user.ID); err == nil {
			isAdmin = adminCheck
		}
	}
	helpText := `*Помощь по командам:*

*/start* - Начать работу с ботом
*/help* - Показать это сообщение
*/subscribe* - Управление подписками на источники
*/news [страница]* - Последние новости из ваших подписок
*/source_news <id> [страница]* - Новости конкретного источника
*/sources* - Все доступные источники новостей
*/categories* - Показать все категории
*/add_source* - Добавить новый источник
*/update* - Обновить новости вручную
*/update_status <id>* - Статус обновления новостей`

	if isAdmin {
		helpText += `*Админские команды:*

• */admin* - Панель администратора
• */admin_users* - Список пользователей
• */admin_stats* - Статистика системы
• */admin_make_admin <id>* - Назначить админа
• */admin_remove_admin <id>* - Снять админа
• */admin_add_category <название>* - Добавить категорию
• */admin_update_source <id> <true/false>* - Изменить активность источника`
	}

	helpText += `

*Горячие кнопки:*
Новости - Последние новости
Мои подписки - Управление подписками
Добавить подписку - Подписаться на новые источники
Обновить новости - Обновить ленту вручную

*Поддержка:*
Если возникли проблемы, напишите @saneknaumchik`

	h.sendMessage(message.Chat.ID, helpText)
}

func (h *Handler) handleSubscribeCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	h.showSubscriptionMenu(ctx, message.Chat.ID, user.ID)
}

func (h *Handler) showUserSubscriptions(ctx context.Context, chatID, userID int64) {
	subscriptions, err := h.service.GetUserSubscriptions(ctx, userID)
	if err != nil {
		h.sendMessage(chatID, "Ошибка получения подписок")
		return
	}

	if len(subscriptions) == 0 {
		h.sendMessage(chatID, "У вас пока нет подписок. Используйте 'Добавить подписку'")
		return
	}

	text := "*Ваши подписки:*\n\n"
	for i, source := range subscriptions {
		text += fmt.Sprintf("%d. %s\n", i+1, source.Name)
		if source.URL != "" {
			text += fmt.Sprintf("%s\n", source.URL)
		}
		text += "\n"
	}

	h.sendMessage(chatID, text)
}

func (h *Handler) handleNewsCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	args := strings.Fields(message.CommandArguments())
	page := 1
	if len(args) > 0 {
		if p, err := strconv.Atoi(args[0]); err == nil && p > 0 {
			page = p
		}
	}

	h.showUserNewsWithPagination(ctx, message.Chat.ID, user.ID, page, 4)
}

func (h *Handler) showUserNewsWithPagination(ctx context.Context, chatID, userID int64, page, pageSize int) {
	response, err := h.service.GetNewsForUserWithPagination(ctx, userID, page, pageSize)
	if err != nil {
		h.sendMessage(chatID, "Ошибка получения новостей")
		return
	}
	if len(response.Data) == 0 {
		if page > 1 {
			h.sendMessage(chatID, "На этой странице больше нет новостей")
		} else {
			h.sendMessage(chatID, "У вас пока нет новостей. Подпишитесь на источники!")
		}
		return
	}

	pageInfo := fmt.Sprintf("*Страница %d из %d*\n\n", response.Page, response.TotalPages)
	h.sendMessage(chatID, pageInfo)

	for i, item := range response.Data {
		text := fmt.Sprintf(
			"*%d. %s*\n\n"+
				"%s (UTC)\n"+
				"%s\n"+
				"[Читать статью](%s)",
			i+1,
			item.Title,
			item.PublishedAt.UTC().Format("02.01.2006 15:04"),
			item.SourceName,
			item.URL,
		)

		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "Markdown"
		// msg.DisableWebPagePreview = true

		if i == len(response.Data)-1 && response.TotalPages > 1 {
			var inlineButtons []tgbotapi.InlineKeyboardButton

			if response.Page > 1 {
				inlineButtons = append(inlineButtons,
					tgbotapi.NewInlineKeyboardButtonData("◀️ Предыдущая",
						fmt.Sprintf("news_page:%d", response.Page-1)))
			}
			if response.Page < response.TotalPages {
				inlineButtons = append(inlineButtons,
					tgbotapi.NewInlineKeyboardButtonData("Следующая ▶️",
						fmt.Sprintf("news_page:%d", response.Page+1)))
			}

			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(inlineButtons...),
			)
			msg.ReplyMarkup = keyboard
		}

		h.bot.Send(msg)
		time.Sleep(100 * time.Millisecond)
	}
}

func (h *Handler) handleSourcesCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	args := strings.Fields(message.CommandArguments())
	page := 1
	if len(args) > 0 {
		if p, err := strconv.Atoi(args[0]); err == nil && p > 0 {
			page = p
		}
	}

	h.showSourcesWithPagination(ctx, message.Chat.ID, page, 10)
}

func (h *Handler) showSourcesWithPagination(ctx context.Context, chatID int64, page, pageSize int) {
	response, err := h.service.GetAllSources(ctx, page, pageSize)
	if err != nil {
		h.sendMessage(chatID, "Ошибка получения источников")
		return
	}

	if len(response.Data) == 0 {
		if page > 1 {
			h.sendMessage(chatID, "На этой странице больше нет источников")
		} else {
			h.sendMessage(chatID, "Источники не найдены")
		}
		return
	}

	text := fmt.Sprintf("*Доступные источники (стр. %d из %d):*\n\n", response.Page, response.TotalPages)

	for _, source := range response.Data {
		status := "✅"
		if !source.IsActive {
			status = "❌"
		}

		text += fmt.Sprintf("%s *ID %d* - %s\n", status, source.ID, source.Name)
		if source.URL != "" {
			text += fmt.Sprintf("  %s\n", source.URL)
		}
		text += fmt.Sprintf("Категория ID: %d\n\n", source.CategoryID)
	}

	if response.TotalPages > 1 {
		text += "\n*Навигация:*\n"

		if response.Page > 1 {
			text += fmt.Sprintf("`/sources %d` - предыдущая страница\n", response.Page-1)
		}

		if response.Page < response.TotalPages {
			text += fmt.Sprintf("`/sources %d` - следующая страница\n", response.Page+1)
		}
	}

	text += "\nℹ️ *Как использовать:*\n"
	text += "• Используйте `/source_news <id>` для просмотра новостей источника\n"
	text += "• Используйте `/add_source Название; URL; ID_категории` для добавления\n"
	text += "• Используйте `/categories` для просмотра всех категорий"

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"

	if response.TotalPages > 1 {
		var inlineButtons []tgbotapi.InlineKeyboardButton

		if response.Page > 1 {
			inlineButtons = append(inlineButtons,
				tgbotapi.NewInlineKeyboardButtonData("◀️ Предыдущая",
					fmt.Sprintf("sources_page:%d", response.Page-1)))
		}

		if response.Page < response.TotalPages {
			inlineButtons = append(inlineButtons,
				tgbotapi.NewInlineKeyboardButtonData("Следующая ▶️",
					fmt.Sprintf("sources_page:%d", response.Page+1)))
		}

		if len(inlineButtons) > 0 {
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(inlineButtons...),
			)
			msg.ReplyMarkup = keyboard
		}
	}

	h.bot.Send(msg)
}

func (h *Handler) handleSourceNewsCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	args := strings.Fields(message.CommandArguments())
	if len(args) == 0 {
		h.sendMessage(message.Chat.ID,
			"Укажите ID источника:\n"+
				"`/source_news <id_источника>`\n\n"+
				"Пример: `/source_news 1`\n"+
				"Используйте `/sources` чтобы посмотреть ID источников")
		return
	}

	sourceID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil || sourceID <= 0 {
		h.sendMessage(message.Chat.ID, "ID источника должно быть положительным числом")
		return
	}

	page := 1
	if len(args) > 1 {
		if p, err := strconv.Atoi(args[1]); err == nil && p > 0 {
			page = p
		}
	}

	h.showSourceNewsWithPagination(ctx, message.Chat.ID, user.ID, sourceID, page, 4)
}

func (h *Handler) showSourceNewsWithPagination(ctx context.Context, chatID, userID, sourceID int64, page, pageSize int) {
	response, err := h.service.GetNewsBySourceWithPagination(ctx, sourceID, userID, page, pageSize)
	if err != nil {
		h.sendMessage(chatID, fmt.Sprintf("%s", err.Error()))
		return
	}

	if len(response.Data) == 0 {
		if page > 1 {
			h.sendMessage(chatID, "На этой странице больше нет новостей")
		} else {
			h.sendMessage(chatID, "У этого источника пока нет новостей")
		}
		return
	}

	source, err := h.service.sourceRepo.GetByID(ctx, int(sourceID))
	if err != nil {
		h.sendMessage(chatID, "Ошибка получения информации об источнике")
		return
	}

	headerText := fmt.Sprintf("*%s*\n", source.Name)
	headerText += fmt.Sprintf("Страница %d из %d\n\n", response.Page, response.TotalPages)

	h.sendMessage(chatID, headerText)

	for i, item := range response.Data {
		text := fmt.Sprintf(
			"*%d. %s*\n\n"+
				"%s (UTC)\n"+
				"[Читать статью](%s)",
			i+1,
			item.Title,
			item.PublishedAt.UTC().Format("02.01.2006 15:04"),
			item.URL,
		)

		msg := tgbotapi.NewMessage(chatID, text)
		msg.ParseMode = "Markdown"
		msg.DisableWebPagePreview = true

		if i == len(response.Data)-1 && response.TotalPages > 1 {
			var inlineButtons []tgbotapi.InlineKeyboardButton

			if response.Page > 1 {
				inlineButtons = append(inlineButtons,
					tgbotapi.NewInlineKeyboardButtonData("◀️ Предыдущая",
						fmt.Sprintf("source_news_nav:%d:%d", sourceID, response.Page-1)))
			}

			inlineButtons = append(inlineButtons,
				tgbotapi.NewInlineKeyboardButtonURL("Открыть статью", item.URL))

			if response.Page < response.TotalPages {
				inlineButtons = append(inlineButtons,
					tgbotapi.NewInlineKeyboardButtonData("Следующая ▶️",
						fmt.Sprintf("source_news_nav:%d:%d", sourceID, response.Page+1)))
			}

			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(inlineButtons...),
			)
			msg.ReplyMarkup = keyboard
		} else {
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL("Открыть статью", item.URL),
				),
			)
			msg.ReplyMarkup = keyboard
		}

		h.bot.Send(msg)
		time.Sleep(100 * time.Millisecond)
	}

	if response.TotalPages > 1 {
		navText := "\n*Навигация по страницам:*\n"
		if response.Page > 1 {
			navText += fmt.Sprintf("`/source_news %d %d` - предыдущая страница\n", sourceID, response.Page-1)
		}
		if response.Page < response.TotalPages {
			navText += fmt.Sprintf("`/source_news %d %d` - следующая страница\n", sourceID, response.Page+1)
		}
		h.sendMessage(chatID, navText)
	}
}

func (h *Handler) handleAddSourceCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	args := message.CommandArguments()
	if args == "" {
		h.sendMessage(message.Chat.ID,
			"Для добавления источника используйте формат:\n"+
				"`/add_source Название; URL; ID_категории`\n\n"+
				"Пример:\n"+
				"`/add_source Habr; https://habr.com/ru/rss/articles/; 1`\n\n"+
				"Для просмотра доступных категорий используйте /categories")
		return
	}

	parts := strings.Split(args, ";")
	if len(parts) != 3 {
		h.sendMessage(message.Chat.ID,
			"Неверный формат. Используйте: Название; URL; ID_категории\n"+
				"Пример: `/add_source Habr; https://habr.com/ru/rss/articles/; 1`")
		return
	}

	name := strings.TrimSpace(parts[0])
	url := strings.TrimSpace(parts[1])
	categoryIDStr := strings.TrimSpace(parts[2])

	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		h.sendMessage(message.Chat.ID, "ID категории должно быть числом")
		return
	}

	err = h.service.AddSource(ctx, name, url, categoryID, user.ID)
	if err != nil {
		h.sendMessage(message.Chat.ID, fmt.Sprintf("Ошибка при добавлении источника: %v", err))
		return
	}

	h.sendMessage(message.Chat.ID, "Источник успешно добавлен!")
}

func (h *Handler) handleCategoriesCommand(ctx context.Context, message *tgbotapi.Message) {
	categories, err := h.service.GetAllCategories(ctx)
	if err != nil {
		h.sendMessage(message.Chat.ID, "Ошибка получения категорий")
		return
	}

	if len(categories) == 0 {
		h.sendMessage(message.Chat.ID, "Категории не найдены")
		return
	}

	text := "*Доступные категории:*\n\n"
	for _, cat := range categories {
		text += fmt.Sprintf("• *ID %d* - %s\n", cat.ID, cat.Name)
	}

	h.sendMessage(message.Chat.ID, text)
}

func (h *Handler) handleUpdateCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	requestID, err := h.service.RequestNewsUpdate(ctx, user.ID)
	if err != nil {
		h.sendMessage(message.Chat.ID, fmt.Sprintf("Ошибка: %v", err))
		return
	}

	h.sendMessage(message.Chat.ID,
		fmt.Sprintf("Запрос на обновление отправлен!\n\n"+
			"ID запроса: `%s`\n"+
			"Используйте `/update_status %s` для проверки статуса",
			requestID, requestID))
}

func (h *Handler) handleUpdateStatusCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	requestID := strings.TrimSpace(message.CommandArguments())
	if requestID == "" {
		h.sendMessage(message.Chat.ID,
			"Укажите ID запроса:\n"+
				"`/update_status <request_id>`\n\n"+
				"ID запроса вы получаете после команды /update")
		return
	}

	req, found := h.service.GetUpdateStatus(ctx, requestID)
	if !found {
		h.sendMessage(message.Chat.ID, "Запрос не найден или устарел")
		return
	}

	var statusText string

	switch req.Status {
	case "pending":
		statusText = "В ожидании"
	case "queued":
		statusText = "В очереди"
	case "processing":
		statusText = "Выполняется"
	case "completed":
		statusText = "Завершено"
	case "failed":
		statusText = "Ошибка"
	default:
		statusText = "Неизвестно"
	}

	text := fmt.Sprintf("📋 *Статус обновления*\n\n"+
		"ID запроса: `%s`\n"+
		"Статус: %s\n"+
		"Пользователь: `%d`\n"+
		"Время запроса: %s\n",
		req.ID, statusText, req.UserID,
		req.Timestamp.Format("02.01.2006 15:04:05"))

	if req.Status == "completed" {
		text += fmt.Sprintf("Результат: %d новых новостей\n", req.Result)
	}

	h.sendMessage(message.Chat.ID, text)
}

func (h *Handler) handleText(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	switch message.Text {
	case "Новости":
		h.showUserNewsWithPagination(ctx, message.Chat.ID, user.ID, 1, 4)
	case "Мои подписки":
		h.showUserSubscriptions(ctx, message.Chat.ID, user.ID)
	case "Добавить подписку":
		h.showSubscriptionMenu(ctx, message.Chat.ID, user.ID)
	case "Обновить новости":
		h.handleUpdateCommand(ctx, message, user)
	case "Источники":
		h.showSourcesWithPagination(ctx, message.Chat.ID, 1, 10)
	case "Помощь":
		h.handleHelp(ctx, message)
	default:
		h.sendMessage(message.Chat.ID, "Я понимаю только команды и кнопки меню. Используйте /help для списка команд.")
	}
}

func (h *Handler) handleCallbackQuery(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	h.bot.Send(tgbotapi.NewCallback(callback.ID, ""))

	data := callback.Data
	chatID := callback.Message.Chat.ID

	user, err := h.service.authService.RegisterOrUpdateTelegramUser(ctx, callback.From.ID,
		callback.From.UserName, callback.From.FirstName+" "+callback.From.LastName)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return
	}

	switch {
	case data == "back_to_main":
		h.showMainMenu(chatID, *user.TgFirstName)

	case strings.HasPrefix(data, "subscribe:"):
		sourceIDStr := strings.TrimPrefix(data, "subscribe:")
		sourceID, err := strconv.Atoi(sourceIDStr)
		if err != nil {
			h.sendMessage(chatID, "Ошибка: неверный ID источника")
			return
		}

		err = h.service.SubscribeUser(ctx, user.ID, sourceID)
		if err != nil {
			h.sendMessage(chatID, fmt.Sprintf("Ошибка подписки: %v", err))
		} else {
			h.sendMessage(chatID, "Подписка оформлена!")
		}

		h.showSubscriptionMenu(ctx, chatID, user.ID)

	case strings.HasPrefix(data, "unsubscribe:"):
		sourceIDStr := strings.TrimPrefix(data, "unsubscribe:")
		sourceID, err := strconv.Atoi(sourceIDStr)
		if err != nil {
			h.sendMessage(chatID, "Ошибка: неверный ID источника")
			return
		}

		err = h.service.UnsubscribeUser(ctx, user.ID, sourceID)
		if err != nil {
			h.sendMessage(chatID, fmt.Sprintf("Ошибка отписки: %v", err))
		} else {
			h.sendMessage(chatID, "Подписка отменена")
		}

		h.showSubscriptionMenu(ctx, chatID, user.ID)

	case strings.HasPrefix(data, "news_page:"):
		pageStr := strings.TrimPrefix(data, "news_page:")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			h.sendMessage(chatID, "Ошибка навигации")
			return
		}
		h.showUserNewsWithPagination(ctx, chatID, user.ID, page, 4)

	case strings.HasPrefix(data, "source_news_nav:"):
		parts := strings.Split(strings.TrimPrefix(data, "source_news_nav:"), ":")
		if len(parts) != 2 {
			h.sendMessage(chatID, "Ошибка навигации")
			return
		}

		sourceID, err1 := strconv.ParseInt(parts[0], 10, 64)
		page, err2 := strconv.Atoi(parts[1])

		if err1 != nil || err2 != nil || sourceID <= 0 || page < 1 {
			h.sendMessage(chatID, "Ошибка навигации")
			return
		}

		h.showSourceNewsWithPagination(ctx, chatID, user.ID, sourceID, page, 4)

	case strings.HasPrefix(data, "sources_page:"):
		pageStr := strings.TrimPrefix(data, "sources_page:")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			h.sendMessage(chatID, "Ошибка навигации")
			return
		}
		h.showSourcesWithPagination(ctx, chatID, page, 10)

	default:
		log.Printf("Неизвестный callback: %s", data)
	}
}

func (h *Handler) showMainMenu(chatID int64, firstName string) {
	text := fmt.Sprintf("*%s*, выберите действие:", firstName)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = MainMenuKeyboard()
	h.bot.Send(msg)
}

func (h *Handler) showSubscriptionMenu(ctx context.Context, chatID, userID int64) {
	sources, err := h.service.GetAllActiveSources(ctx)
	if err != nil {
		h.sendMessage(chatID, "Ошибка получения источников")
		return
	}

	subscriptions, err := h.service.GetUserSubscriptions(ctx, userID)
	if err != nil {
		h.sendMessage(chatID, "Ошибка получения подписок")
		return
	}

	var subscribedIDs []int
	for _, sub := range subscriptions {
		subscribedIDs = append(subscribedIDs, int(sub.ID))
	}

	keyboard := SubscriptionKeyboard(sources, subscribedIDs)

	text := "*Управление подписками*\n\n" +
		"Нажмите на источник чтобы изменить подписку:"
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	h.bot.Send(msg)
}
