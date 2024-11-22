package model

import (
	"fmt"
	"strings"
)

type TonAddr struct {
	Hex                  string
	MainnetBounceable    string
	MainnetNonBounceale  string
	TestnetBounceable    string
	TestnetNonBounceable string
}

func (t TonAddr) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*Hex*: %s\n", EscapeMarkdownV2(t.Hex)))
	sb.WriteString(fmt.Sprintf("*MainnetBounceable*: %s\n", EscapeMarkdownV2(t.MainnetBounceable)))
	sb.WriteString(fmt.Sprintf("*MainnetNonBounceale*: `%s`\n", EscapeMarkdownV2(t.MainnetNonBounceale)))
	sb.WriteString(fmt.Sprintf("*TestnetBounceable*: `%s`\n", EscapeMarkdownV2(t.TestnetBounceable)))
	sb.WriteString(fmt.Sprintf("*TestnetNonBounceable*: `%s`\n", EscapeMarkdownV2(t.TestnetNonBounceable)))
	return sb.String()
}
