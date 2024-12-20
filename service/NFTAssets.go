package service

import (
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ttv-bot/model"
)

func (s *Service) handleNFTAssetsCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "💠 Please enter the *wallet address*, as shown in this example:\n" + model.EscapeMarkdownV2("EQAVKjQ5Gkj2GvO0_QpvlW-iC1WhiIG428zywVcP3-C7LPEi") + "\n"
	msg.ParseMode = api.ModeMarkdown
	s.SendMessage(msg)
}

func (s *Service) replyNFTAssetsCommand(upmsg *api.Message) {
	s.replyNFTAssetsPage(upmsg, 1, true)
}

func (s *Service) replyNFTAssetsPage(upmsg *api.Message, page int, first bool) {
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
					api.NewInlineKeyboardButtonData("+1 ➡️", fmt.Sprintf("%s-%d", model.NFTAssetCommand, page+1+1+1)),
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
		nfts := s.nftItemsCache.get(reMsg, addr, maxPage*pageSize, s.queryNFTAssetsByWallet)
		if nfts == nil || len(nfts) == 0 {
			text = "💠	The wallet doesn't hold any NFTs"
		} else {
			text = "We found " + fmt.Sprintf("%+v", nfts[0].Count) + "💎 on your request:\n"
			totalPage := len(nfts) / pageSize
			if len(nfts)%pageSize != 0 {
				totalPage++
			}
			if page >= totalPage {
				keyboard = api.NewInlineKeyboardMarkup(
					api.NewInlineKeyboardRow(
						api.NewInlineKeyboardButtonData("-1 ⬅️", fmt.Sprintf("%s-%d", model.NFTAssetCommand, page-1)),
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
					api.NewInlineKeyboardButtonData("-1 ⬅️", fmt.Sprintf("%s-%d", model.NFTAssetCommand, page+1-1)),
					api.NewInlineKeyboardButtonData(fmt.Sprintf("%+v", page+1), "current"),
					api.NewInlineKeyboardButtonData("+1 ➡️", fmt.Sprintf("%s-%d", model.NFTAssetCommand, page+1+1)),
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
