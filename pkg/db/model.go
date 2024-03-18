package db

import "time"

type Link struct {
	ShortSuffix      string     `json:"short_suffix"`
	Link             string     `json:"link"`
	SecretKey        string     `json:"secret_key"`
	ShortLinkExpDate *time.Time `json:"expiration_date"`
	Clicks           int        `json:"clicks"`
}
