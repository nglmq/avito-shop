package models

type InventoryItem struct {
	Item     string `db:"item" json:"type"`
	Quantity int    `db:"quantity" json:"quantity"`
}
