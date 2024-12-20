package service

import (
	"encoding/json"
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ttv-bot/model"
)

func (s *Service) handleFakeCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "⁉️ Please select the type you want to query for fake\n"
	keyboardMarkup := api.NewInlineKeyboardMarkup(api.NewInlineKeyboardRow(api.NewInlineKeyboardButtonData("👻Account", model.FakeAccountCommand), api.NewInlineKeyboardButtonData("💰Jetton", model.FakeJettonCommand), api.NewInlineKeyboardButtonData("💠NFT", model.FakeNFTCommand)))
	msg.ReplyMarkup = keyboardMarkup
	s.SendMessage(msg)
}

func (s *Service) handleFakeAccountCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "👻 Please enter the *Account address* and *name* as shown in this example:\n\n----------\n`{\"address\":\"EQCMOXxD-f8LSWWbXQowKxqTr3zMY-X1wMTyWp3B-LR6s3Va\",\"name\":\"Telegram\"}`"
	msg.ParseMode = api.ModeMarkdown
	s.SendMessage(msg)
}

func (s *Service) handleFakeNFTCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "👻 Please enter the *NFT address* and *name* as shown in this example:\n\n----------\n`{\"address\":\"EQAZKtHcN6mhbzhbzcnuj-94r5P-hZYBKEwQ4_-dc-2AWMqZ\",\"name\":\"Spinners\"}`"
	msg.ParseMode = api.ModeMarkdown
	s.SendMessage(msg)
}

func (s *Service) handleFakeJettonCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "👻 Please enter the *Jetton address* and *name* as shown in this example:\n\n----------\n`{\"address\":\"EQCMOXxD-f8LSWWbXQowKxqTr3zMY-X1wMTyWp3B-LR6s3Va\",\"name\":\"STON\"}`"
	msg.ParseMode = api.ModeMarkdown
	s.SendMessage(msg)
}

func (s *Service) replyFakeAccountCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.ReplyToMessageID = upmsg.MessageID
	var acc model.Account
	err := json.Unmarshal([]byte(upmsg.Text), &acc)
	if err != nil {
		msg.Text = "invalid input. " + err.Error()
	} else {
		err, tonAddr := model.ParseTonAddress(acc.Address)
		if err != nil {
			msg.Text = "invalid address: " + model.EscapeMarkdownV2(acc.Address) + ", " + err.Error()
		} else {
			prMsg := api.NewMessage(upmsg.Chat.ID, "🔍 processing....")
			prMsg.ReplyToMessageID = upmsg.MessageID
			s.SendMessage(prMsg)
			if name := s.validAccountsCache.GetAccountNameByAddress(tonAddr); name != "" {
				if name == acc.Name {
					msg.Text = "⭕️ *Valid*\n\n*Name*:" + acc.Name + "\n" + fmt.Sprintf("%+v", tonAddr) + "\n"
				} else {
					msg.Text = "❗️*Fake*\n\n*Name*:" + name + "\n" + fmt.Sprintf("%+v", tonAddr) + "\n"
					if addr := s.validAccountsCache.GetAccountAddressByName(acc.Name); addr != nil {
						msg.Text = msg.Text + "\n*Name*:" + acc.Name + "\n" + fmt.Sprintf("%+v", addr) + "\n"
					}
				}
			} else if addr := s.validAccountsCache.GetAccountAddressByName(acc.Name); addr != nil {
				if addr.Hex != tonAddr.Hex {
					msg.Text = "❌ *Fake*\n\n*Name*:" + acc.Name + "\n" + fmt.Sprintf("%+v", addr) + "\n"
				}
			} else {
				msg.Text = "❓Unable to verify account\n" + fmt.Sprintf("%+v", acc) + "\n"
			}
		}

	}
	msg.ParseMode = api.ModeMarkdown
	s.SendMessage(msg)
}

func (s *Service) replyFakeJettonCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	var acc model.Account
	err := json.Unmarshal([]byte(upmsg.Text), &acc)
	if err != nil {
		msg.Text = "invalid input. " + err.Error()
	} else {
		err, tonAddr := model.ParseTonAddress(acc.Address)
		if err != nil {
			msg.Text = "invalid address: " + model.EscapeMarkdownV2(acc.Address) + ", " + err.Error()
		} else {
			prMsg := api.NewMessage(upmsg.Chat.ID, "🔍 processing....")
			prMsg.ReplyToMessageID = upmsg.MessageID
			s.SendMessage(prMsg)
			if name := s.validJettonsCache.GetJettonNameByAddress(tonAddr); name != "" {
				if name == acc.Name {
					msg.Text = "⭕️ *Valid*\n\n*Name*:" + acc.Name + "\n" + fmt.Sprintf("%+v", tonAddr) + "\n"
				} else {
					msg.Text = "❗️*Fake*\n\n*Name*:" + name + "\n" + fmt.Sprintf("%+v", tonAddr) + "\n"
					if addr := s.validJettonsCache.GetJettonAddressByName(acc.Name); addr != nil {
						msg.Text = msg.Text + "\n*Name*:" + acc.Name + "\n" + fmt.Sprintf("%+v", addr) + "\n"
					}
				}
			} else if addr := s.validJettonsCache.GetJettonAddressByName(acc.Name); addr != nil {
				if addr.Hex != tonAddr.Hex {
					msg.Text = "❌ *Fake*\n\n*Name*:" + acc.Name + "\n" + fmt.Sprintf("%+v", addr) + "\n"
				}
			} else {
				msg.Text = "❓Unable to verify jetton\n" + fmt.Sprintf("%+v", acc) + "\n"
			}
		}

	}
	msg.ParseMode = api.ModeMarkdown
	s.SendMessage(msg)
}

func (s *Service) replyFakeNFTCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	var acc model.Account
	err := json.Unmarshal([]byte(upmsg.Text), &acc)
	if err != nil {
		msg.Text = "invalid input. " + err.Error()
	} else {
		err, tonAddr := model.ParseTonAddress(acc.Address)
		if err != nil {
			msg.Text = "invalid address: " + model.EscapeMarkdownV2(acc.Address) + ", " + err.Error()
		} else {
			prMsg := api.NewMessage(upmsg.Chat.ID, "🔍 processing....")
			prMsg.ReplyToMessageID = upmsg.MessageID
			s.SendMessage(prMsg)
			if name := s.validNFTsCache.GetNFTNameByAddress(tonAddr); name != "" {
				if name == acc.Name {
					msg.Text = "⭕️ *Valid*\n\n*Name*:" + acc.Name + "\n" + fmt.Sprintf("%+v", tonAddr) + "\n"
				} else {
					msg.Text = "❗️*Fake*\n\n*Name*:" + name + "\n" + fmt.Sprintf("%+v", tonAddr) + "\n"
					if addr := s.validNFTsCache.GetNFTAddressByName(acc.Name); addr != nil {
						msg.Text = msg.Text + "\n*Name*:" + acc.Name + "\n" + fmt.Sprintf("%+v", addr) + "\n"
					}
				}
			} else if addr := s.validNFTsCache.GetNFTAddressByName(acc.Name); addr != nil {
				if addr.Hex != tonAddr.Hex {
					msg.Text = "❌ *Fake*\n\n*Name*:" + acc.Name + "\n" + fmt.Sprintf("%+v", addr) + "\n"
				}
			} else {
				msg.Text = "❓Unable to verify nft\n" + fmt.Sprintf("%+v", acc) + "\n"
			}
		}

	}
	msg.ParseMode = api.ModeMarkdown
	s.SendMessage(msg)
}
