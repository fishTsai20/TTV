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
	keyboardMarkup := api.NewInlineKeyboardMarkup(api.NewInlineKeyboardRow(api.NewInlineKeyboardButtonData("👻Account", "fakeAccount"), api.NewInlineKeyboardButtonData("💰Jetton", "fakeJetton"), api.NewInlineKeyboardButtonData("💠NFT", "fakeNFT")))
	msg.ReplyMarkup = keyboardMarkup
	s.SendMessage(msg)
}

func (s *Service) handleFakeTypeCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "👻 Please enter the *address* and *name* as shown in this example:\n\n----------\n`{\"address\":\"EQCMOXxD-f8LSWWbXQowKxqTr3zMY-X1wMTyWp3B-LR6s3Va\",\"name\":\"Telegram\"}`"
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
			if name := s.validCache.GetAccountNameByAddress(tonAddr); name != "" {
				if name == acc.Name {
					msg.Text = "⭕️ *Valid*\n\n*Name*:" + acc.Name + "\n" + fmt.Sprintf("%+v", tonAddr) + "\n"
				} else {
					msg.Text = "❗️*Fake*\n\n*Name*:" + name + "\n" + fmt.Sprintf("%+v", tonAddr) + "\n"
					if addr := s.validCache.GetAccountAddressByName(acc.Name); addr != nil {
						msg.Text = msg.Text + "\n*Name*:" + acc.Name + "\n" + fmt.Sprintf("%+v", addr) + "\n"
					}
				}
			} else if addr := s.validCache.GetAccountAddressByName(acc.Name); addr != nil {
				if addr.Hex != tonAddr.Hex {
					msg.Text = "❌ *Fake*\n\n*Name*:" + acc.Name + "\n" + fmt.Sprintf("%+v", addr) + "\n"
				}
			} else {
				msg.Text = "❓Unable to verify \n" + fmt.Sprintf("%+v", acc) + "\n"
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
			if name := s.validCache.GetJettonNameByAddress(tonAddr); name != "" {
				if name == acc.Name {
					msg.Text = "⭕️ *Valid*\n\n*Name*:" + acc.Name + "\n" + fmt.Sprintf("%+v", tonAddr) + "\n"
				} else {
					msg.Text = "❗️*Fake*\n\n*Name*:" + name + "\n" + fmt.Sprintf("%+v", tonAddr) + "\n"
					if addr := s.validCache.GetJettonAddressByName(acc.Name); addr != nil {
						msg.Text = msg.Text + "\n*Name*:" + acc.Name + "\n" + fmt.Sprintf("%+v", addr) + "\n"
					}
				}
			} else if addr := s.validCache.GetJettonAddressByName(acc.Name); addr != nil {
				if addr.Hex != tonAddr.Hex {
					msg.Text = "❌ *Fake*\n\n*Name*:" + acc.Name + "\n" + fmt.Sprintf("%+v", addr) + "\n"
				}
			} else {
				msg.Text = "❓Unable to verify \n" + fmt.Sprintf("%+v", acc) + "\n"
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
			if name := s.validCache.GetNFTNameByAddress(tonAddr); name != "" {
				if name == acc.Name {
					msg.Text = "⭕️ *Valid*\n\n*Name*:" + acc.Name + "\n" + fmt.Sprintf("%+v", tonAddr) + "\n"
				} else {
					msg.Text = "❗️*Fake*\n\n*Name*:" + name + "\n" + fmt.Sprintf("%+v", tonAddr) + "\n"
					if addr := s.validCache.GetNFTAddressByName(acc.Name); addr != nil {
						msg.Text = msg.Text + "\n*Name*:" + acc.Name + "\n" + fmt.Sprintf("%+v", addr) + "\n"
					}
				}
			} else if addr := s.validCache.GetNFTAddressByName(acc.Name); addr != nil {
				if addr.Hex != tonAddr.Hex {
					msg.Text = "❌ *Fake*\n\n*Name*:" + acc.Name + "\n" + fmt.Sprintf("%+v", addr) + "\n"
				}
			} else {
				msg.Text = "❓Unable to verify \n" + fmt.Sprintf("%+v", acc) + "\n"
			}
		}

	}
	msg.ParseMode = api.ModeMarkdown
	s.SendMessage(msg)
}
