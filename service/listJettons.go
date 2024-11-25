package service

import (
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ttv-bot/model"
)

/**
1. listJettons
*/
// call back command /listjettons
func (s *Service) handleListJettonsCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "💰 Do you want to list Jettons base on name ?\n"
	keyboardMarkup := api.NewInlineKeyboardMarkup(api.NewInlineKeyboardRow(api.NewInlineKeyboardButtonData("Y", model.ListJettonsByNameCommand), api.NewInlineKeyboardButtonData("N", model.ListJettonsDefaultCommand)))
	msg.ReplyMarkup = keyboardMarkup
	s.SendMessage(msg)
}

/**
2-1-1 listJettons with name filter
*/
// callback command /listJettonsByName -> receive name and reply
func (s *Service) handleListJettonsByNameCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "💰 Please enter Jetton name\n"
	s.SendMessage(msg)
}

/**
2-1-2 listJettons with name filter---- first
*/
// process message reply to command /listJettonsByName
func (s *Service) replyListJettonsByNameCommand(upmsg *api.Message) {
	s.replyListJettonsPage(upmsg, 1, true)
}

/**
2-2-1 listJettons without any filters
*/
// callback command /listJettonsDefault -> page
func (s *Service) handleListJettonsDefaultCommand(upmsg *api.Message) {
	s.replyListJettonsPageWithName(upmsg, 1, true, "")
}

/**
3 callback: reply to page
*/
// callback
func (s *Service) replyListJettonsPage(upmsg *api.Message, page int, first bool) {
	var name string
	if first {
		name = upmsg.Text
	} else {
		if upmsg.ReplyToMessage != nil {
			name = upmsg.ReplyToMessage.Text
		} else {
			name = ""
		}
	}
	s.replyListJettonsPageWithName(upmsg, page, first, name)
}

func (s *Service) replyListJettonsPageWithName(upmsg *api.Message, page int, first bool, name string) {
	var (
		msg      api.Chattable
		keyboard api.InlineKeyboardMarkup
	)

	page--
	if page < 0 {
		keyboard = api.NewInlineKeyboardMarkup(
			api.NewInlineKeyboardRow(
				api.NewInlineKeyboardButtonData(fmt.Sprintf("%+v", page+1+1), "current"),
				api.NewInlineKeyboardButtonData("+1 ➡️", fmt.Sprintf("%s-%d", model.ListJettonsCommand, page+1+1+1)),
			),
		)
		msgk := api.NewEditMessageReplyMarkup(upmsg.Chat.ID, upmsg.MessageID, keyboard)
		s.SendMessage(msgk)
		return
	}
	pageSize := 3
	text := ""
	var jettons map[string]*model.TonAddr
	if name != "" {
		jettons = s.validJettonsCache.FuzzyGET(name)
	} else {
		jettons = s.validJettonsCache.GetAll()
	}
	if jettons == nil || len(jettons) == 0 {
		text = "💰	There is no Jetton to list"
		s.chatMap.Delete(model.ListJettonsCommand)
	} else {
		totalPage := len(jettons) / pageSize
		if len(jettons)%pageSize != 0 {
			totalPage++
		}
		if page >= totalPage {
			keyboard = api.NewInlineKeyboardMarkup(
				api.NewInlineKeyboardRow(
					api.NewInlineKeyboardButtonData("-1 ⬅️", fmt.Sprintf("%s-%d", model.ListJettonsCommand, page-1)),
					api.NewInlineKeyboardButtonData(fmt.Sprintf("%+v", page), "current"),
				),
			)
			msgk := api.NewEditMessageReplyMarkup(upmsg.Chat.ID, upmsg.MessageID, keyboard)
			s.SendMessage(msgk)
			return
		}
		var accs []model.Account
		for name, addr := range jettons {
			accs = append(accs, model.Account{
				Address: addr.MainnetBounceable,
				Name:    name,
			})
		}
		text += "We found " + fmt.Sprintf("%+v", len(accs)) + " jettons on your request:\n"
		text += s.formatPage(model.ConvertToTgTextSlice(accs), page, pageSize)
		keyboard = api.NewInlineKeyboardMarkup(
			api.NewInlineKeyboardRow(
				api.NewInlineKeyboardButtonData("-1 ⬅️", fmt.Sprintf("%s-%d", model.ListJettonsCommand, page+1-1)),
				api.NewInlineKeyboardButtonData(fmt.Sprintf("%+v", page+1), "current"),
				api.NewInlineKeyboardButtonData("+1 ➡️", fmt.Sprintf("%s-%d", model.ListJettonsCommand, page+1+1)),
			),
		)

	}
	if !first {
		editMsg := api.NewEditMessageText(upmsg.Chat.ID, upmsg.MessageID, "")
		if len(keyboard.InlineKeyboard) != 0 {
			editMsg.ReplyMarkup = &keyboard
			editMsg.DisableWebPagePreview = true
		}
		editMsg.ParseMode = api.ModeMarkdown
		editMsg.Text = text
		msg = editMsg
	} else {
		replyMsg := api.NewMessage(upmsg.Chat.ID, "")
		if len(keyboard.InlineKeyboard) != 0 {
			replyMsg.ReplyMarkup = keyboard
			replyMsg.DisableWebPagePreview = true
		}
		replyMsg.ParseMode = api.ModeMarkdown
		//list jettons with name filter
		if name != "" {
			replyMsg.ReplyToMessageID = upmsg.MessageID
		}
		replyMsg.Text = text
		msg = replyMsg
	}

	s.SendMessage(msg)
}
