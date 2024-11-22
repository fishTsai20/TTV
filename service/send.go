package service

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

func (s *Service) SendMessage(msg api.Chattable) api.Message {
	switch v := msg.(type) {
	case api.MessageConfig:
		if v.Text == "" {
			return api.Message{}
		}
		v.Text = v.Text
		mmsg, err := s.bot.Send(v)
		if err != nil {
			log.Println(err)
		}
		return mmsg
	case api.EditMessageTextConfig:
		if v.Text == "" {
			return api.Message{}
		}
		v.Text = v.Text
		mmsg, err := s.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		return mmsg
	case api.EditMessageReplyMarkupConfig:
		if v.ReplyMarkup == nil {
			return api.Message{}
		}
		mmsg, err := s.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		return mmsg
	}
	return api.Message{}
}

func (s *Service) deleteMessage(gid int64, mid int) {
	time.Sleep(time.Second * 240)
	deleteMessageConfig := api.DeleteMessageConfig{
		ChatID:    gid,
		MessageID: mid,
	}
	_, err := s.bot.Request(deleteMessageConfig)

	if err != nil {
		log.Println(err.Error())
	}
}
