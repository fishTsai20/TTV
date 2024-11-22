package service

import (
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ttv-bot/model"
)

func (s *Service) handle24HJettonNewPoolsCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.ReplyToMessageID = upmsg.MessageID
	pools := s.stonFiPoolsCache.get(msg, s.query24HJettonNewPools)
	if pools == nil || len(pools) == 0 {
		msg.Text = "🏊 There is no new pools created by ston.fi within 24h"
		s.chatMap.Delete(model.JetttonNewPoolsCommand)
	} else {
		msg.Text = "We found " + fmt.Sprintf("%+v", len(pools)) + "🏊 on your request:\n"
		msg.Text += s.formatPage(model.ConvertToTgTextSlice(pools), 0, 3)
		keyboard := api.NewInlineKeyboardMarkup(
			api.NewInlineKeyboardRow(
				api.NewInlineKeyboardButtonData("1", "current"),
				api.NewInlineKeyboardButtonData("+1 ➡️", fmt.Sprintf("%s-%d", model.JetttonNewPoolsCommand, 2)),
			),
		)
		msg.ReplyMarkup = keyboard
		msg.ParseMode = api.ModeMarkdown
		msg.DisableWebPagePreview = true
	}
	s.SendMessage(msg)
}

func (s *Service) reply24HJettonNewPoolsCommand(upmsg *api.Message, page int, first bool) {
	page--
	if page < 0 {
		keyboard := api.NewInlineKeyboardMarkup(
			api.NewInlineKeyboardRow(
				api.NewInlineKeyboardButtonData(fmt.Sprintf("%+v", page+1+1), "current"),
				api.NewInlineKeyboardButtonData("+1 ➡️", fmt.Sprintf("%s-%d", model.JetttonNewPoolsCommand, page+1+1+1)),
			),
		)
		msgk := api.NewEditMessageReplyMarkup(upmsg.Chat.ID, upmsg.MessageID, keyboard)
		s.SendMessage(msgk)
		return
	}
	pageSize := 3
	msg := api.NewEditMessageText(upmsg.Chat.ID, upmsg.MessageID, "")
	pools := s.stonFiPoolsCache.get(msg, s.query24HJettonNewPools)
	if pools == nil || len(pools) == 0 {
		msg.Text = "🏊 There is no new pools created by ston.fi within 24h"
		s.chatMap.Delete(model.JetttonNewPoolsCommand)
	} else {
		totalPage := len(pools) / pageSize
		if len(pools)%pageSize != 0 {
			totalPage++
		}
		if page >= totalPage {
			keyboard := api.NewInlineKeyboardMarkup(
				api.NewInlineKeyboardRow(
					api.NewInlineKeyboardButtonData("-1 ⬅️", fmt.Sprintf("%s-%d", model.JetttonNewPoolsCommand, page-1)),
					api.NewInlineKeyboardButtonData(fmt.Sprintf("%+v", page), "current"),
				),
			)
			msgk := api.NewEditMessageReplyMarkup(upmsg.Chat.ID, upmsg.MessageID, keyboard)
			s.SendMessage(msgk)
			return
		}
		msg.Text += s.formatPage(model.ConvertToTgTextSlice(pools), page, pageSize)
		keyboard := api.NewInlineKeyboardMarkup(
			api.NewInlineKeyboardRow(
				api.NewInlineKeyboardButtonData("-1 ⬅️", fmt.Sprintf("%s-%d", model.JetttonNewPoolsCommand, page+1-1)),
				api.NewInlineKeyboardButtonData(fmt.Sprintf("%+v", page+1), "current"),
				api.NewInlineKeyboardButtonData("+1 ➡️", fmt.Sprintf("%s-%d", model.JetttonNewPoolsCommand, page+1+1)),
			),
		)
		msg.ReplyMarkup = &keyboard
		msg.ParseMode = api.ModeMarkdown
		msg.DisableWebPagePreview = true
	}
	s.SendMessage(msg)
}
