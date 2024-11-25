package service

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ttv-bot/model"
)

func (b *Service) handleStartCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "🤖 Hey hey, " + model.EscapeMarkdownV2(getUserName(upmsg.From)) + ", *Welcome to Ton Token Vigil\\(TTV\\)*\\!\n"
	msg.Text += "📣 Comprehensive data analysis\\.  \n📈 Real\\-time wallet P&L analysis\\.  \n⚠️ Detect insider trading/rug risks\\.  \n\n*Powered by:*  \n⛽️ [Chainbase](https://chainbase.com/)"
	keyboardMarkup := api.NewInlineKeyboardMarkup(api.NewInlineKeyboardRow(api.NewInlineKeyboardButtonData("🐝	honeyPot", model.HoneyPotCommand), api.NewInlineKeyboardButtonData("⁉️	fake", model.FakeCommand)), api.NewInlineKeyboardRow(api.NewInlineKeyboardButtonData("📖	list", model.ListCommand), api.NewInlineKeyboardButtonData("💰	jettons", model.JettonsCommand), api.NewInlineKeyboardButtonData("💠	nfts", model.NFTCommand)))
	msg.ReplyMarkup = keyboardMarkup
	msg.ParseMode = api.ModeMarkdownV2
	b.SendMessage(msg)
	b.chatMap.Delete(upmsg.Chat.ID)
}
