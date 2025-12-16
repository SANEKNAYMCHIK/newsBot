package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleAdminCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	isAdmin, err := h.service.IsAdmin(ctx, user.ID)
	if err != nil || !isAdmin {
		h.sendMessage(message.Chat.ID, "Эта команда только для администраторов")
		return
	}

	text := `*Панель администратора*

*Основные команды:*
• */admin_users [страница]* - Список пользователей
• */admin_stats* - Статистика системы
• */admin_make_admin <user_id>* - Назначить админа
• */admin_remove_admin <user_id>* - Снять админа

*Управление категориями:*
• */admin_add_category <название>* - Добавить категорию
• */categories* - Показать все категории

*Управление источниками:*
• */admin_update_source <id> <true/false>* - Изменить активность источника
• */sources* - Показать все источники

*Мониторинг:*
• */update_status <request_id>* - Проверить статус обновления

Для помощи по конкретной команде используйте */help*`

	h.sendMessage(message.Chat.ID, text)
}

func (h *Handler) handleAdminUsersCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	isAdmin, err := h.service.IsAdmin(ctx, user.ID)
	if err != nil || !isAdmin {
		h.sendMessage(message.Chat.ID, "Эта команда только для администраторов")
		return
	}

	args := strings.Fields(message.CommandArguments())
	page := 1
	pageSize := 20

	if len(args) >= 1 {
		if p, err := strconv.Atoi(args[0]); err == nil {
			page = p
		}
	}
	if len(args) >= 2 {
		if ps, err := strconv.Atoi(args[1]); err == nil {
			pageSize = ps
		}
	}

	users, err := h.service.GetUsers(ctx, page, pageSize)
	if err != nil {
		h.sendMessage(message.Chat.ID, fmt.Sprintf("Ошибка получения пользователей: %v", err))
		return
	}

	if len(users.Data) == 0 {
		h.sendMessage(message.Chat.ID, "Пользователи не найдены")
		return
	}

	text := fmt.Sprintf("*Список пользователей (стр. %d):*\n\n", page)
	for i, u := range users.Data {
		email := "N/A"
		if u.Email != nil {
			email = *u.Email
		}
		tgUsername := "N/A"
		if u.TgUsername != nil {
			tgUsername = *u.TgUsername
		}

		text += fmt.Sprintf("*%d.*ID: `%d`\n", i+1, u.ID)
		text += fmt.Sprintf("Email: `%s`\n", email)
		text += fmt.Sprintf("TG: @%s\n", tgUsername)
		text += fmt.Sprintf("Роль: `%s`\n\n", u.Role)
	}

	text += fmt.Sprintf("*Всего:* %d пользователей\n", users.Total)
	text += fmt.Sprintf("*Страниц:* %d", users.TotalPages)

	h.sendMessage(message.Chat.ID, text)
}

func (h *Handler) handleAdminMakeAdminCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	isAdmin, err := h.service.IsAdmin(ctx, user.ID)
	if err != nil || !isAdmin {
		h.sendMessage(message.Chat.ID, "Эта команда только для администраторов")
		return
	}

	targetUserID, err := strconv.ParseInt(message.CommandArguments(), 10, 64)
	if err != nil || targetUserID == 0 {
		h.sendMessage(message.Chat.ID,
			"Неверный формат. Используйте:\n"+
				"`/admin_make_admin <user_id>`\n\n"+
				"Пример: `/admin_make_admin 123456`")
		return
	}

	err = h.service.MakeAdmin(ctx, targetUserID, user.ID)
	if err != nil {
		h.sendMessage(message.Chat.ID, fmt.Sprintf("Ошибка: %v", err))
		return
	}

	h.sendMessage(message.Chat.ID, fmt.Sprintf("Пользователь с ID `%d` назначен администратором", targetUserID))
}

func (h *Handler) handleAdminRemoveAdminCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	isAdmin, err := h.service.IsAdmin(ctx, user.ID)
	if err != nil || !isAdmin {
		h.sendMessage(message.Chat.ID, "Эта команда только для администраторов")
		return
	}

	targetUserID, err := strconv.ParseInt(message.CommandArguments(), 10, 64)
	if err != nil || targetUserID == 0 {
		h.sendMessage(message.Chat.ID,
			"Неверный формат. Используйте:\n"+
				"`/admin_remove_admin <user_id>`\n\n"+
				"Пример: `/admin_remove_admin 123456`")
		return
	}

	err = h.service.RemoveAdmin(ctx, targetUserID, user.ID)
	if err != nil {
		h.sendMessage(message.Chat.ID, fmt.Sprintf("Ошибка: %v", err))
		return
	}

	h.sendMessage(message.Chat.ID, fmt.Sprintf("У пользователя с ID `%d` сняты права администратора", targetUserID))
}

func (h *Handler) handleAdminAddCategoryCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	isAdmin, err := h.service.IsAdmin(ctx, user.ID)
	if err != nil || !isAdmin {
		h.sendMessage(message.Chat.ID, "Эта команда только для администраторов")
		return
	}

	categoryName := strings.TrimSpace(message.CommandArguments())
	if categoryName == "" {
		h.sendMessage(message.Chat.ID,
			"Укажите название категории:\n"+
				"`/admin_add_category <название>`\n\n"+
				"Пример: `/admin_add_category Технологии`")
		return
	}

	category, err := h.service.CreateCategory(ctx, categoryName)
	if err != nil {
		h.sendMessage(message.Chat.ID, fmt.Sprintf("Ошибка создания категории: %v", err))
		return
	}

	h.sendMessage(message.Chat.ID,
		fmt.Sprintf("Категория создана:\n"+
			"ID: `%d`\n"+
			"Название: *%s*", category.ID, category.Name))
}

func (h *Handler) handleAdminUpdateSourceCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	isAdmin, err := h.service.IsAdmin(ctx, user.ID)
	if err != nil || !isAdmin {
		h.sendMessage(message.Chat.ID, "Эта команда только для администраторов")
		return
	}

	args := strings.Fields(message.CommandArguments())
	if len(args) != 2 {
		h.sendMessage(message.Chat.ID,
			"Неверный формат. Используйте:\n"+
				"`/admin_update_source <id_источника> <true/false>`\n\n"+
				"Примеры:\n"+
				"`/admin_update_source 1 true` - активировать источник\n"+
				"`/admin_update_source 1 false` - деактивировать источник")
		return
	}

	sourceID, err := strconv.Atoi(args[0])
	if err != nil {
		h.sendMessage(message.Chat.ID, "ID источника должно быть числом")
		return
	}

	isActive := strings.ToLower(args[1])
	if isActive != "true" && isActive != "false" {
		h.sendMessage(message.Chat.ID, "Второй параметр должен быть true или false")
		return
	}

	activeBool := isActive == "true"
	err = h.service.UpdateSource(ctx, sourceID, activeBool)
	if err != nil {
		h.sendMessage(message.Chat.ID, fmt.Sprintf("Ошибка обновления источника: %v", err))
		return
	}

	status := "активирован"
	if !activeBool {
		status = "деактивирован"
	}

	h.sendMessage(message.Chat.ID,
		fmt.Sprintf("Источник с ID `%d` успешно %s", sourceID, status))
}

func (h *Handler) handleAdminStatsCommand(ctx context.Context, message *tgbotapi.Message, user *models.User) {
	isAdmin, err := h.service.IsAdmin(ctx, user.ID)
	if err != nil || !isAdmin {
		h.sendMessage(message.Chat.ID, "Эта команда только для администраторов")
		return
	}

	stats, err := h.service.GetSystemStats(ctx)
	if err != nil {
		h.sendMessage(message.Chat.ID, fmt.Sprintf("Ошибка получения статистики: %v", err))
		return
	}

	text := "📊 *Статистика системы*\n\n"
	text += fmt.Sprintf("Пользователей: *%d*\n", stats["users_count"])
	text += fmt.Sprintf("Источников: *%d*\n", stats["sources_count"])
	text += fmt.Sprintf("Новостей: *%d*\n", stats["news_count"])

	h.sendMessage(message.Chat.ID, text)
}
