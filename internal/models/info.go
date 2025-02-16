package models

type InfoResponse struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}

type CoinHistory struct {
	Received []TransactionReceivedHistory `json:"received"`
	Sent     []TransactionSentHistory     `json:"sent"`
}

type TransactionReceivedHistory struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type TransactionSentHistory struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}
