package service

import (
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ttv-bot/model"
)

/**
1. listNFTs
*/
// call back command /listnfts
func (s *Service) handleListNFTsCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "💰 Do you want to list NFTs base on name ?\n"
	keyboardMarkup := api.NewInlineKeyboardMarkup(api.NewInlineKeyboardRow(api.NewInlineKeyboardButtonData("Y", model.ListNFTsByNameCommand), api.NewInlineKeyboardButtonData("N", model.ListNFTsDefaultCommand)))
	msg.ReplyMarkup = keyboardMarkup
	s.SendMessage(msg)
}

/**
2-1-1 listNFTs with name filter
*/
// callback command /listNFTsByName -> receive name and reply
func (s *Service) handleListNFTsByNameCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "💰 Please enter NFT name\n"
	s.SendMessage(msg)
}

/**
2-1-2 listNFTs with name filter---- first
*/
// process message reply to command /listNFTsByName
func (s *Service) replyListNFTsByNameCommand(upmsg *api.Message) {
	s.replyListNFTsPage(upmsg, 1, true)
}

/**
2-2-1 listNFTs without any filters
*/
// callback command /listNFTsDefault -> page
func (s *Service) handleListNFTsDefaultCommand(upmsg *api.Message) {
	s.replyListNFTsPageWithName(upmsg, 1, true, "")
}

/**
3 callback: reply to page
*/
// callback
func (s *Service) replyListNFTsPage(upmsg *api.Message, page int, first bool) {
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
	s.replyListNFTsPageWithName(upmsg, page, first, name)
}

func (s *Service) replyListNFTsPageWithName(upmsg *api.Message, page int, first bool, name string) {
	var (
		msg      api.Chattable
		keyboard api.InlineKeyboardMarkup
	)

	page--
	if page < 0 {
		keyboard = api.NewInlineKeyboardMarkup(
			api.NewInlineKeyboardRow(
				api.NewInlineKeyboardButtonData(fmt.Sprintf("%+v", page+1+1), "current"),
				api.NewInlineKeyboardButtonData("+1 ➡️", fmt.Sprintf("%s-%d", model.ListNFTsCommand, page+1+1+1)),
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
		jettons = s.validNFTsCache.FuzzyGET(name)
	} else {
		jettons = s.validNFTsCache.GetAll()
	}
	if jettons == nil || len(jettons) == 0 {
		text = "💰	There is no NFT to list"
		s.chatMap.Delete(model.ListNFTsCommand)
	} else {
		totalPage := len(jettons) / pageSize
		if len(jettons)%pageSize != 0 {
			totalPage++
		}
		if page >= totalPage {
			keyboard = api.NewInlineKeyboardMarkup(
				api.NewInlineKeyboardRow(
					api.NewInlineKeyboardButtonData("-1 ⬅️", fmt.Sprintf("%s-%d", model.ListNFTsCommand, page-1)),
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
		text += "We found " + fmt.Sprintf("%+v", len(accs)) + " NFTs on your request:\n"
		text += s.formatPage(model.ConvertToTgTextSlice(accs), page, pageSize)
		keyboard = api.NewInlineKeyboardMarkup(
			api.NewInlineKeyboardRow(
				api.NewInlineKeyboardButtonData("-1 ⬅️", fmt.Sprintf("%s-%d", model.ListNFTsCommand, page+1-1)),
				api.NewInlineKeyboardButtonData(fmt.Sprintf("%+v", page+1), "current"),
				api.NewInlineKeyboardButtonData("+1 ➡️", fmt.Sprintf("%s-%d", model.ListNFTsCommand, page+1+1)),
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
		//list nfts with name filter
		if name != "" {
			replyMsg.ReplyToMessageID = upmsg.MessageID
		}
		replyMsg.Text = text
		msg = replyMsg
	}

	s.SendMessage(msg)
}
