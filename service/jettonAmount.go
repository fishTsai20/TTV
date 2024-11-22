package service

import (
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ttv-bot/model"
)

func (s *Service) handleJettonAmountCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "🛰️	Please enter the *jetton master address*\n"
	msg.ParseMode = api.ModeMarkdown
	s.SendMessage(msg)
}

func (s *Service) replyJettonAmountCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.ReplyToMessageID = upmsg.MessageID
	msg.ParseMode = api.ModeMarkdown
	contract := upmsg.Text
	err, addr := model.ParseTonAddress(contract)
	if err != nil {
		msg.Text = "invalid address: " + contract + ", " + err.Error() + "\n"
	} else {
		vol, err := s.queryJetton24hVol(addr, msg)
		if err != nil {
			msg.Text = "failed to get jetton's change\n"
		} else {
			msg.Text = "🛰️ *24h amount*: " + model.FormatNumber(vol) + "\n\n" + fmt.Sprintf("%+v", addr) + "\n"
		}
	}
	s.SendMessage(msg)
}
