package service

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ttv-bot/model"
)

func (s *Service) handleListCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "📖 Please select the type you want to list\n"
	keyboardMarkup := api.NewInlineKeyboardMarkup(api.NewInlineKeyboardRow(api.NewInlineKeyboardButtonData("👻Account", model.ListAccountsCommand), api.NewInlineKeyboardButtonData("💰Jetton", model.ListJettonsCommand), api.NewInlineKeyboardButtonData("💠NFT", model.ListNFTsCommand)))
	msg.ReplyMarkup = keyboardMarkup
	s.SendMessage(msg)
}
