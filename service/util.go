package service

import (
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"ttv-bot/model"
)

func getUserName(user *api.User) string {
	if user != nil {
		if user.UserName != "" {
			return user.UserName
		}
		return user.FirstName
	}
	return ""
}

func (s *Service) formatPage(datas []model.TgText, page int, pageSize int) string {
	start := page * pageSize
	end := start + pageSize
	if end > len(datas) {
		end = len(datas)
	}
	datas = datas[start:end]
	res := ""
	for _, pool := range datas {
		res += pool.ToTgText()
	}
	return res
}
