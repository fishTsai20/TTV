package service

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"ttv-bot/model"
)

func (s *Service) handleHoneyPotCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "🐝️	Please enter the *account address*\n"
	msg.ParseMode = api.ModeMarkdown
	s.SendMessage(msg)
}

func (s *Service) replyHoneyPotCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.ReplyToMessageID = upmsg.MessageID
	msg.ParseMode = api.ModeMarkdown
	contract := upmsg.Text
	err, addr := model.ParseTonAddress(contract)
	if err != nil {
		msg.Text = "invalid address: " + contract + ", " + err.Error() + "\n"
	} else {
		deployed, err := s.queryIsContractDeployed(addr, msg)
		if err != nil {
			msg.Text = "failed to judge account deployed or not\n"
		} else {
			if deployed {
				msg.Text = "️ 🟢 deployed \n"
			} else {
				msg.Text = "️ ⛔️ undeployed \n"
			}
		}
		msg.Text += "️ 💡 todo: is contract code verified \n"
		msg.Text += "️ 💡 todo: E2E tests \n"
	}
	s.SendMessage(msg)
}
