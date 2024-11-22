package service

import (
	"encoding/json"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"ttv-bot/log"
	"ttv-bot/model"
)

func (s *Service) handleJettonBalanceCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "🐟 Please enter the *jetton_master_address* and *wallet_address* as shown in this example:\n\n---------\n`{\"jetton_master_address\":\"EQB5Wjo7yXdaB70yBoN2YEv8iVPjAdMObf_Dq40ELLaPllNb\",\"wallet_address\":\"EQCP4Qoe6kstey5LlQCQ6uocuhwBhJ2ylFSgipy4fqQHOlsN\"}`"
	msg.ParseMode = api.ModeMarkdown
	s.SendMessage(msg)
}

func (s *Service) replyJettonBalanceCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.ReplyToMessageID = upmsg.MessageID
	msg.ParseMode = api.ModeMarkdown
	var jetton model.Jetton
	err := json.Unmarshal([]byte(upmsg.Text), &jetton)
	if err != nil {
		msg.Text = "invalid input. " + err.Error()
	} else {
		err, contract := model.ParseTonAddress(jetton.JettonMasterAddress)
		if err != nil {
			msg.Text = "invalid jettonMasterAddress: " + model.EscapeMarkdownV2(jetton.JettonMasterAddress) + ", " + err.Error()
		} else {
			reMsg := api.NewMessage(upmsg.Chat.ID, "")
			reMsg.ReplyToMessageID = upmsg.MessageID
			reMsg.ParseMode = api.ModeMarkdown
			err, wallet := model.ParseTonAddress(jetton.WalletAddress)
			if err != nil {
				msg.Text = "invalid walletAddress: " + model.EscapeMarkdownV2(jetton.WalletAddress) + ", " + err.Error()
			} else {
				balance, err := s.queryJettonBalance(reMsg, contract, wallet)
				if err != nil {
					log.Info("failed to query jetton balance", zap.Error(err))
				} else if balance == nil {
					msg.Text = "🐟 Can't get holder balance"
				} else {
					msg.DisableWebPagePreview = true
					msg.Text += balance.ToTgText()
				}
			}
		}
	}
	s.SendMessage(msg)
}
