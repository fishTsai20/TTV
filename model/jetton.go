package model

type Jetton struct {
	JettonWalletAddress string `json:"jetton_wallet_address"`
	WalletAddress       string `json:"wallet_address"`
	JettonMasterAddress string `json:"jetton_master_address"`
	JettonBalance       string `json:"jetton_balance"`
	Percent             string `json:"percent"`
	Rnn                 int    `json:"rnn"`
}

func (j Jetton) ToTgText() string {
	res := "\n"
	res += "🐳\n*jetton_wallet_address: *[" + j.JettonWalletAddress + "](https://tonscan.org/address/" + j.JettonWalletAddress + ")\n"
	res += "*wallet_address: *[" + j.WalletAddress + "](https://tonscan.org/address/" + j.WalletAddress + ")\n"
	res += "*jetton_balance: *" + j.JettonBalance + "\n"
	if j.Percent != "" {
		res += "*percent: *" + j.Percent + "%\n"
	}
	return res
}
