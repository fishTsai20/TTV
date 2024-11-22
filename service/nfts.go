package service

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ttv-bot/model"
)

func (s *Service) handleNFTsCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "💠 Please select the type you want to query for nfts\n"
	keyboardMarkup := api.NewInlineKeyboardMarkup(api.NewInlineKeyboardRow(api.NewInlineKeyboardButtonData("💠	collection", model.NFTCollectionCommand), api.NewInlineKeyboardButtonData("💎	nft", model.NFTItemCommand), api.NewInlineKeyboardButtonData("👛	assets", model.NFTAssetCommand)))
	msg.ReplyMarkup = keyboardMarkup
	s.SendMessage(msg)
}
