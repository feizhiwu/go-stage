package model

type Pages struct {
	Count uint `json:"count"`
	Limit uint `json:"limit"`
	Page  uint `json:"page"`
}
