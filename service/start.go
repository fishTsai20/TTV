package service

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"ttv-bot/model"
)

func (b *Service) handleStartCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "🤖 Hey hey, " + getUserName(upmsg.From) + ", welcome on board!\n"
	keyboardMarkup := api.NewInlineKeyboardMarkup(api.NewInlineKeyboardRow(api.NewInlineKeyboardButtonData("🐝	honeyPot", model.HoneyPotCommand), api.NewInlineKeyboardButtonData("⁉️	fake", model.FakeCommand)), api.NewInlineKeyboardRow(api.NewInlineKeyboardButtonData("💰	jettons", model.JettonsCommand), api.NewInlineKeyboardButtonData("💠	nfts", model.NFTCommand)))
	msg.ReplyMarkup = keyboardMarkup
	b.SendMessage(msg)
	b.chatMap.Delete(upmsg.Chat.ID)
}
