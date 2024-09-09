package models

import "time"

type Secret struct {
	ID         int       `json:"id"`
	Identifier string    `json:"identifier"`
	TextSecret string    `json:"text_secret"`
	Counter    uint8     `json:"counter"`
	SaveDate   time.Time `json:"save_time"`
}
