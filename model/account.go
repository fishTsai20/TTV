package model

import (
	"fmt"
	"github.com/xssnick/tonutils-go/address"
	"strings"
)

type Account struct {
	Address string `yaml:"address" json:"address"`
	Name    string `yaml:"name" json:"name"`
}

func ParseTonAddress(addr string) (error, *TonAddr) {
	addr = strings.TrimSpace(addr)
	var (
		parseAddr     *address.Address
		hex           string
		mainBounce    string
		mainNonBounce string
		testBounce    string
		testNonBounce string
		err           error
	)
	if strings.HasPrefix(addr, "0:") || strings.HasPrefix(addr, "1:") {
		parseAddr, err = address.ParseRawAddr(addr)
		if err != nil {
			return err, nil
		}
	} else {
		parseAddr, err = address.ParseAddr(addr)
		if err != nil {
			return err, nil
		}
	}

	if parseAddr.IsTestnetOnly() {
		if parseAddr.IsBounceable() {
			//test-bounceable
			testBounce = parseAddr.String()
			//test-non-bounceable
			parseAddr.SetBounce(false)
			testNonBounce = parseAddr.String()
			//main-non-bouceable
			parseAddr.SetTestnetOnly(false)
			mainNonBounce = parseAddr.String()
			//main-bounceable
			parseAddr.SetBounce(true)
			mainBounce = parseAddr.String()
		} else {
			//test-non-bounceable
			testNonBounce = parseAddr.String()
			//test-bounceable
			parseAddr.SetBounce(true)
			testBounce = parseAddr.String()
			//main-bounceable
			parseAddr.SetTestnetOnly(false)
			mainBounce = parseAddr.String()
			//main-non-bouceable
			parseAddr.SetBounce(false)
			mainNonBounce = parseAddr.String()
		}
	} else {
		if parseAddr.IsBounceable() {
			//main-bounceable
			mainBounce = parseAddr.String()
			//main-non-bouceable
			parseAddr.SetBounce(false)
			mainNonBounce = parseAddr.String()
			//test-non-bounceable
			parseAddr.SetTestnetOnly(true)
			testNonBounce = parseAddr.String()
			//test-bounceable
			parseAddr.SetBounce(true)
			testBounce = parseAddr.String()
		} else {
			//main-non-bouceable
			mainNonBounce = parseAddr.String()
			//main-bounceable
			parseAddr.SetBounce(true)
			mainBounce = parseAddr.String()
			//test-bounceable
			parseAddr.SetTestnetOnly(true)
			testBounce = parseAddr.String()
			//test-non-bounceable
			parseAddr.SetBounce(false)
			testNonBounce = parseAddr.String()
		}

	}
	hex = fmt.Sprintf("%d:%x", int8(parseAddr.Workchain()), parseAddr.Data())
	return nil, &TonAddr{
		Hex:                  hex,
		MainnetBounceable:    mainBounce,
		MainnetNonBounceale:  mainNonBounce,
		TestnetBounceable:    testBounce,
		TestnetNonBounceable: testNonBounce,
	}

}

func (a Account) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*Address*: %s\n", EscapeMarkdownV2(a.Address)))
	sb.WriteString(fmt.Sprintf("*Name*: %s\n", a.Name))
	return sb.String()
}

func (a Account) ToTgText() string {
	res := "\n"
	res += "*Name: *" + a.Name + "\n"
	res += "*Address: *[" + a.Address + "](https://tonscan.org/address/" + a.Address + ")\n"
	return res
}
