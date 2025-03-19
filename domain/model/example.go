package model

import (
	"time"
)

// Example represents a basic example entity
type Example struct {
	Id        int       `json:"id" gorm:"column:id;primaryKey;autoIncrement" structs:",omitempty,underline"`
	Name      string    `json:"name" gorm:"column:name;type:varchar(255);not null" structs:",omitempty,underline"`
	Alias     string    `json:"alias" gorm:"column:alias;type:varchar(255)" structs:",omitempty,underline"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime" structs:",omitempty,underline"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime" structs:",omitempty,underline"`
}

// TableName returns the table name for the Example model
func (e Example) TableName() string {
	return "example"
}
