package entity

import (
    "time"

    "gorm.io/gorm"
)

/**
 * @author Rancho
 * @date 2021/12/24
 */

type Example struct {
    Id        int                    `gorm:"primarykey" structs:",omitempty,underline"`
    Name      string                 `json:"name" structs:",omitempty,underline"`
    Alias     string                 `json:"alias" structs:",omitempty,underline"`
    CreatedAt time.Time              `json:"created_at" structs:",omitempty,underline"`
    UpdatedAt time.Time              `json:"updated_at" structs:",omitempty,underline"`
    DeletedAt gorm.DeletedAt         `json:"deleted_at" structs:",omitempty,underline"`
    ChangeMap map[string]interface{} `gorm:"-" json:"-" structs:"-"`
}

func (e Example) TableName() string {
    return "example"
}
