package entity

import "time"

type ShortUrlMapping struct {
	Id         int64     `json:"id"`
	ShortUrl   string    `json:"shortUrl"`
	LongUrl    string    `json:"longUrl"`
	ExpiryDate time.Time `json:"expiryDate"`
	CreatedAt  time.Time `json:"createdAt"`
}
