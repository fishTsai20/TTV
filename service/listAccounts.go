package service

import (
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ttv-bot/model"
)

/**
1. listAccounts
*/
// call back command /listaccounts
func (s *Service) handleListAccountsCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "💰 Do you want to list Accounts base on name ?\n"
	keyboardMarkup := api.NewInlineKeyboardMarkup(api.NewInlineKeyboardRow(api.NewInlineKeyboardButtonData("Y", model.ListAccountsByNameCommand), api.NewInlineKeyboardButtonData("N", model.ListAccountsDefaultCommand)))
	msg.ReplyMarkup = keyboardMarkup
	s.SendMessage(msg)
}

/**
2-1-1 listAccounts with name filter
*/
// callback command /listAccountsByName -> receive name and reply
func (s *Service) handleListAccountsByNameCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "💰 Please enter Account name\n"
	s.SendMessage(msg)
}

/**
2-1-2 listAccounts with name filter---- first
*/
// process message reply to command /listAccountsByName
func (s *Service) replyListAccountsByNameCommand(upmsg *api.Message) {
	s.replyListAccountsPage(upmsg, 1, true)
}

/**
2-2-1 listAccounts without any filters
*/
// callback command /listAccountsDefault -> page
func (s *Service) handleListAccountsDefaultCommand(upmsg *api.Message) {
	s.replyListAccountsPageWithName(upmsg, 1, true, "")
}

/**
3 callback: reply to page
*/
// callback
func (s *Service) replyListAccountsPage(upmsg *api.Message, page int, first bool) {
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
	s.replyListAccountsPageWithName(upmsg, page, first, name)
}

func (s *Service) replyListAccountsPageWithName(upmsg *api.Message, page int, first bool, name string) {
	var (
		msg      api.Chattable
		keyboard api.InlineKeyboardMarkup
	)

	page--
	if page < 0 {
		keyboard = api.NewInlineKeyboardMarkup(
			api.NewInlineKeyboardRow(
				api.NewInlineKeyboardButtonData(fmt.Sprintf("%+v", page+1+1), "current"),
				api.NewInlineKeyboardButtonData("+1 ➡️", fmt.Sprintf("%s-%d", model.ListAccountsCommand, page+1+1+1)),
			),
		)
		msgk := api.NewEditMessageReplyMarkup(upmsg.Chat.ID, upmsg.MessageID, keyboard)
		s.SendMessage(msgk)
		return
	}
	pageSize := 3
	text := ""
	var accounts map[string]*model.TonAddr
	if name != "" {
		accounts = s.validAccountsCache.FuzzyGET(name)
	} else {
		accounts = s.validAccountsCache.GetAll()
	}
	if accounts == nil || len(accounts) == 0 {
		text = "💰	There is no Accounts to list"
		s.chatMap.Delete(model.ListAccountsCommand)
	} else {
		totalPage := len(accounts) / pageSize
		if len(accounts)%pageSize != 0 {
			totalPage++
		}
		if page >= totalPage {
			keyboard = api.NewInlineKeyboardMarkup(
				api.NewInlineKeyboardRow(
					api.NewInlineKeyboardButtonData("-1 ⬅️", fmt.Sprintf("%s-%d", model.ListAccountsCommand, page-1)),
					api.NewInlineKeyboardButtonData(fmt.Sprintf("%+v", page), "current"),
				),
			)
			msgk := api.NewEditMessageReplyMarkup(upmsg.Chat.ID, upmsg.MessageID, keyboard)
			s.SendMessage(msgk)
			return
		}
		var accs []model.Account
		for name, addr := range accounts {
			accs = append(accs, model.Account{
				Address: addr.MainnetBounceable,
				Name:    name,
			})
		}
		text += "We found " + fmt.Sprintf("%+v", len(accs)) + " accounts on your request:\n"
		text += s.formatPage(model.ConvertToTgTextSlice(accs), page, pageSize)
		keyboard = api.NewInlineKeyboardMarkup(
			api.NewInlineKeyboardRow(
				api.NewInlineKeyboardButtonData("-1 ⬅️", fmt.Sprintf("%s-%d", model.ListAccountsCommand, page+1-1)),
				api.NewInlineKeyboardButtonData(fmt.Sprintf("%+v", page+1), "current"),
				api.NewInlineKeyboardButtonData("+1 ➡️", fmt.Sprintf("%s-%d", model.ListAccountsCommand, page+1+1)),
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
		//list accounts with name filter
		if name != "" {
			replyMsg.ReplyToMessageID = upmsg.MessageID
		}
		replyMsg.Text = text
		msg = replyMsg
	}

	s.SendMessage(msg)
}
