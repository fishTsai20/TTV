package model

import (
	"fmt"
	"strings"
)

func FormatNumber(value float64) string {
	if value < 1000 {
		return fmt.Sprintf("%.2f", value)
	} else if value < 1000000 {
		return fmt.Sprintf("%.2fk", value/1000)
	} else if value < 1000000000 {
		return fmt.Sprintf("%.2fm", value/1000000)
	} else {
		return fmt.Sprintf("%.2fb", value/1000000000)
	}
}

func EscapeMarkdownV2(input string) string {
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "=", "|", "{", "}", ".", "!"}
	for _, char := range specialChars {
		input = strings.ReplaceAll(input, char, "\\"+char)
	}
	return input
}

type TgText interface {
	ToTgText() string
}

func ConvertToTgTextSlice[T TgText](items []T) []TgText {
	result := make([]TgText, len(items))
	for i, item := range items {
		result[i] = item
	}
	return result
}
