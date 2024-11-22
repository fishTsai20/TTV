package service

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"ttv-bot/model"
)

func (s *Service) handleJettonsCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "💰 Please select the type you want to query for jettons\n"
	keyboardMarkup := api.NewInlineKeyboardMarkup(api.NewInlineKeyboardRow(api.NewInlineKeyboardButtonData("🏊	24h stonfi.pools", model.JetttonNewPoolsCommand), api.NewInlineKeyboardButtonData("🤝	24h change", model.JettonChangesCommand), api.NewInlineKeyboardButtonData("🛰️	24h amount", model.JettonAmountCommand)),
		api.NewInlineKeyboardRow(api.NewInlineKeyboardButtonData("👥	holders", model.JettonHolderCommand), api.NewInlineKeyboardButtonData("🐳	Top10", model.JettonTopHoldersCommand), api.NewInlineKeyboardButtonData("🐟	balance", model.JettonBalanceCommand)))
	msg.ReplyMarkup = keyboardMarkup
	s.SendMessage(msg)
}
