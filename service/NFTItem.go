package service

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api"
	"ttv-bot/model"
)

func (s *Service) handleNFTItemCommand(upmsg *api.Message) {
	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.Text = "💎 Please enter the *nft item address*\n"
	msg.ParseMode = api.ModeMarkdown
	s.SendMessage(msg)
}

func (s *Service) replyNFTItemCommand(upmsg *api.Message) {
	contract := upmsg.Text
	err, addr := model.ParseTonAddress(contract)
	text := ""
	if err != nil {
		text = "invalid address: " + contract + ", " + err.Error() + "\n"
	} else {
		reMsg := api.NewMessage(upmsg.Chat.ID, "")
		reMsg.ReplyToMessageID = upmsg.MessageID
		reMsg.ParseMode = api.ModeMarkdown
		nft := s.nftCollectionCache.get(reMsg, addr, s.queryNFTCollectionByItem)
		if nft == nil {
			text = "💎	Failed to get NFT's collection"
		} else {
			text = nft.ToTgText()
		}
	}

	msg := api.NewMessage(upmsg.Chat.ID, "")
	msg.ReplyToMessageID = upmsg.MessageID
	msg.DisableWebPagePreview = true
	msg.ParseMode = api.ModeMarkdown
	msg.Text = text
	s.SendMessage(msg)
}
