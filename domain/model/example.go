package model

import (
	"time"
)

type Example struct {
	Id        int       `json:"id" structs:",omitempty,underline"`
	Name      string    `json:"name" structs:",omitempty,underline"`
	Alias     string    `json:"alias" structs:",omitempty,underline"`
	CreatedAt time.Time `json:"created_at" structs:",omitempty,underline"`
	UpdatedAt time.Time `json:"updated_at" structs:",omitempty,underline"`
}

func (e Example) TableName() string {
	return "example"
}
