package model

type Pool struct {
	TxHash    string `json:"tx_hash"`
	Pool      string `json:"pool"`
	CreatedAt string `json:"createdAt"`
}

func (p Pool) ToTgText() string {
	res := "\n"
	res += "🏊	*pool: *[" + p.Pool + "](https://tonscan.org/jetton/" + p.Pool + ")\n"
	res += "----------\n*tx_hash: *[" + p.TxHash + "](https://tonscan.org/tx/" + p.TxHash + ")\n"
	res += "*created_at: *" + p.CreatedAt + "\n"
	return res
}
