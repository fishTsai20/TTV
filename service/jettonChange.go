package service

import (
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"ttv-bot/log"
	"ttv-bot/model"
)

func (s *Service) handleJettonChangeCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "🤝	Please enter the *jetton master address*as shown in this example, as shown in this example:\n" + model.EscapeMarkdownV2("EQBynBO23ywHy_CgarY9NK9FTz0yDsG82PtcbSTQgGoXwiuA")
	msg.ParseMode = api.ModeMarkdown
	s.SendMessage(msg)
}

func (s *Service) replyJettonChangeCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.ParseMode = api.ModeMarkdown
	msg.ReplyToMessageID = upmsg.MessageID
	contract := upmsg.Text
	err, addr := model.ParseTonAddress(contract)
	if err != nil {
		msg.Text = "invalid address, " + err.Error() + "\n"
	} else {
		totalAmount, err := s.queryJettonTotalAmount(addr, msg)
		if err != nil || totalAmount == -1 {
			log.Info("query total amount failed", zap.Error(err))
			msg.Text = "failed to get jetton's change\n"
		} else {
			vol, err := s.queryJetton24hVol(addr, msg)
			if err != nil {
				msg.Text = "failed to get jetton's change\n"
			} else {
				change := vol * 2 / totalAmount
				msg.Text = "🤝 *24h change*: " + model.FormatNumber(change) + "%\n\n" + fmt.Sprintf("%+v", addr) + "\n"
			}
		}
	}

	s.SendMessage(msg)
}
