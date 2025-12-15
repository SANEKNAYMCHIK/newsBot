package bot

import (
	"fmt"

	"github.com/SANEKNAYMCHIK/newsBot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Новости"),
			tgbotapi.NewKeyboardButton("Мои подписки"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить подписку"),
			tgbotapi.NewKeyboardButton("Обновить новости"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Источники"),
			tgbotapi.NewKeyboardButton("Помощь"),
		),
	)
}

func RemoveKeyboard() tgbotapi.ReplyKeyboardRemove {
	return tgbotapi.ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      false,
	}
}

func SubscriptionKeyboard(sources []models.Source, subscribedSources []int) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, source := range sources {
		isSubscribed := false
		for _, subID := range subscribedSources {
			if int(source.ID) == subID {
				isSubscribed = true
				break
			}
		}

		var buttonText string
		var callbackData string

		if isSubscribed {
			buttonText = fmt.Sprintf("✅ %s", source.Name)
			callbackData = fmt.Sprintf("unsubscribe:%d", source.ID)
		} else {
			buttonText = fmt.Sprintf("❌ %s", source.Name)
			callbackData = fmt.Sprintf("subscribe:%d", source.ID)
		}

		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(buttonText, callbackData),
		)
		rows = append(rows, row)
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "back_to_main"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func CategoryKeyboard(categories []models.Category) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, category := range categories {
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(category.Name, fmt.Sprintf("category:%d", category.ID)),
		)
		rows = append(rows, row)
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", "back_to_main"),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func NewsNavigationKeyboard(newsID int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Читать статью", ""),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Предыдущая", "prev_news"),
			tgbotapi.NewInlineKeyboardButtonData("Следующая", "next_news"),
		),
	)
}
