package service

import (
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ttv-bot/model"
)

func (s *Service) handleNFTCollectionCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "💠 Please enter the *nft collection address*\n"
	msg.ParseMode = api.ModeMarkdown
	s.SendMessage(msg)
}

func (s *Service) replyNFTCollectionsCommand(upmsg *api.Message) {
	s.replyNFTCollectionPage(upmsg, 1, true)
}

func (s *Service) replyNFTCollectionPage(upmsg *api.Message, page int, first bool) {
	var (
		msg      api.Chattable
		keyboard api.InlineKeyboardMarkup
		contract string
	)
	if !first {
		contract = upmsg.ReplyToMessage.Text
	} else {
		contract = upmsg.Text
	}
	err, addr := model.ParseTonAddress(contract)
	text := ""
	if err != nil {
		text = "invalid address, " + err.Error() + "\n"
	} else {
		page--
		if page < 0 {
			keyboard = api.NewInlineKeyboardMarkup(
				api.NewInlineKeyboardRow(
					api.NewInlineKeyboardButtonData(fmt.Sprintf("%+v", page+1+1), "current"),
					api.NewInlineKeyboardButtonData("+1 ➡️", fmt.Sprintf("%s-%d", model.NFTCollectionCommand, page+1+1+1)),
				),
			)
			msgk := api.NewEditMessageReplyMarkup(upmsg.Chat.ID, upmsg.MessageID, keyboard)
			s.SendMessage(msgk)
			return
		}
		pageSize := 3
		maxPage := 10
		reMsg := api.NewMessage(upmsg.Chat.ID, "")
		reMsg.ReplyToMessageID = upmsg.MessageID
		nfts := s.nftItemsCache.get(reMsg, addr, maxPage*pageSize, s.queryNFTitemsByCollection)
		if nfts == nil || len(nfts) == 0 {
			text = "💠	There is no NFT items"
		} else {
			text = "We found " + fmt.Sprintf("%+v", nfts[0].Count) + "💎 on your request:\n"
			totalPage := len(nfts) / pageSize
			if len(nfts)%pageSize != 0 {
				totalPage++
			}
			if page >= totalPage {
				keyboard = api.NewInlineKeyboardMarkup(
					api.NewInlineKeyboardRow(
						api.NewInlineKeyboardButtonData("-1 ⬅️", fmt.Sprintf("%s-%d", model.NFTCollectionCommand, page-1)),
						api.NewInlineKeyboardButtonData(fmt.Sprintf("%+v", page), "current"),
					),
				)
				msgk := api.NewEditMessageReplyMarkup(upmsg.Chat.ID, upmsg.MessageID, keyboard)
				s.SendMessage(msgk)
				return
			}
			text += s.formatPage(model.ConvertToTgTextSlice(nfts), page, pageSize)
			keyboard = api.NewInlineKeyboardMarkup(
				api.NewInlineKeyboardRow(
					api.NewInlineKeyboardButtonData("-1 ⬅️", fmt.Sprintf("%s-%d", model.NFTCollectionCommand, page+1-1)),
					api.NewInlineKeyboardButtonData(fmt.Sprintf("%+v", page+1), "current"),
					api.NewInlineKeyboardButtonData("+1 ➡️", fmt.Sprintf("%s-%d", model.NFTCollectionCommand, page+1+1)),
				),
			)

		}
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
		replyMsg.ReplyToMessageID = upmsg.MessageID
		replyMsg.Text = text
		msg = replyMsg
	}

	s.SendMessage(msg)
}
