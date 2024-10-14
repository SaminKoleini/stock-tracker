package models

type Stock struct {
	Symbol string `json:"symbol"`
	Price  int    `json:"price"`
	Volume int    `json:"volume"`
}
