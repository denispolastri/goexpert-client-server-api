package types

import "gorm.io/gorm"

type Dollar struct {
	gorm.Model
	Code      string  `json:"code"`
	CodeIn    string  `json:"codeIn"`
	Name      string  `json:"name"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	VarBid    float64 `json:"varBid"`
	PctChange float64 `json:"pctChange"`
	Bid       float64 `json:"bid"`
	Ask       float64 `json:"ask"`
	Timestamp string  `json:"timestamp"`
}
