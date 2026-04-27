package link

import "time"

type Link struct {
	Hash string `json:"hash"`
	Url  string `josn:"url"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
