package service

import (
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ttv-bot/model"
)

func (s *Service) handleJettonHoldersCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "👥 Please enter the *jetton master address*as shown in this example, as shown in this example:\n" + model.EscapeMarkdownV2("EQBynBO23ywHy_CgarY9NK9FTz0yDsG82PtcbSTQgGoXwiuA")
	msg.ParseMode = api.ModeMarkdown
	s.SendMessage(msg)
}

func (s *Service) replyJettonHoldersCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.ReplyToMessageID = upmsg.MessageID
	contract := upmsg.Text
	err, addr := model.ParseTonAddress(contract)
	if err != nil {
		msg.Text = "invalid address, " + err.Error() + "\n"
	} else {
		holder, err := s.queryJettonHolderCountAmount(addr, msg)
		if err != nil || holder == -1 {
			msg.Text = "failed to get jetton's holders\n"
		} else {
			msg.Text = "👥 *Holders*: " + fmt.Sprintf("%+v", holder) + "\n\n" + fmt.Sprintf("%+v", addr) + "\n"
		}
	}
	msg.ParseMode = api.ModeMarkdown
	s.SendMessage(msg)
}
