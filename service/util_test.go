package service

import (
	"github.com/stretchr/testify/require"
	"testing"
	"ttv-bot/model"
)

func TestTestnetNonBounceable(t *testing.T) {
	_, tonAddr := model.ParseTonAddress("0QC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nMsH")
	require.Equal(t, tonAddr.TestnetNonBounceable, "0QC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nMsH")
	require.Equal(t, tonAddr.Hex, "0:ba295e33b3c4c9b5265aa4ead1166a92931ce9abea120a8c5e91044a1257f89c")
	require.Equal(t, tonAddr.TestnetBounceable, "kQC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nJbC")
	require.Equal(t, tonAddr.MainnetNonBounceale, "UQC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nHCN")
	require.Equal(t, tonAddr.MainnetBounceable, "EQC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nC1I")
}

func TestTestnetBounceable(t *testing.T) {
	_, tonAddr := model.ParseTonAddress("kQC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nJbC")
	require.Equal(t, tonAddr.TestnetNonBounceable, "0QC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nMsH")
	require.Equal(t, tonAddr.Hex, "0:ba295e33b3c4c9b5265aa4ead1166a92931ce9abea120a8c5e91044a1257f89c")
	require.Equal(t, tonAddr.TestnetBounceable, "kQC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nJbC")
	require.Equal(t, tonAddr.MainnetNonBounceale, "UQC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nHCN")
	require.Equal(t, tonAddr.MainnetBounceable, "EQC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nC1I")
}

func TestMainnetBounceable(t *testing.T) {
	_, tonAddr := model.ParseTonAddress("EQC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nC1I")
	require.Equal(t, tonAddr.TestnetNonBounceable, "0QC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nMsH")
	require.Equal(t, tonAddr.Hex, "0:ba295e33b3c4c9b5265aa4ead1166a92931ce9abea120a8c5e91044a1257f89c")
	require.Equal(t, tonAddr.TestnetBounceable, "kQC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nJbC")
	require.Equal(t, tonAddr.MainnetNonBounceale, "UQC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nHCN")
	require.Equal(t, tonAddr.MainnetBounceable, "EQC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nC1I")
}

func TestMainnetNonBounceable(t *testing.T) {
	_, tonAddr := model.ParseTonAddress("UQC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nHCN")
	require.Equal(t, tonAddr.TestnetNonBounceable, "0QC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nMsH")
	require.Equal(t, tonAddr.Hex, "0:ba295e33b3c4c9b5265aa4ead1166a92931ce9abea120a8c5e91044a1257f89c")
	require.Equal(t, tonAddr.TestnetBounceable, "kQC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nJbC")
	require.Equal(t, tonAddr.MainnetNonBounceale, "UQC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nHCN")
	require.Equal(t, tonAddr.MainnetBounceable, "EQC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nC1I")
}

func TestHex(t *testing.T) {
	_, tonAddr := model.ParseTonAddress("0:ba295e33b3c4c9b5265aa4ead1166a92931ce9abea120a8c5e91044a1257f89c")
	require.Equal(t, tonAddr.TestnetNonBounceable, "0QC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nMsH")
	require.Equal(t, tonAddr.Hex, "0:ba295e33b3c4c9b5265aa4ead1166a92931ce9abea120a8c5e91044a1257f89c")
	require.Equal(t, tonAddr.TestnetBounceable, "kQC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nJbC")
	require.Equal(t, tonAddr.MainnetNonBounceale, "UQC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nHCN")
	require.Equal(t, tonAddr.MainnetBounceable, "EQC6KV4zs8TJtSZapOrRFmqSkxzpq-oSCoxekQRKElf4nC1I")
}
