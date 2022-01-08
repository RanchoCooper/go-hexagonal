package entity

import (
    "gorm.io/gorm"
)

/**
 * @author Rancho
 * @date 2021/12/24
 */

type Example struct {
    gorm.Model
    Name      string                 `json:"name" structs:",underline"`
    Alias     string                 `json:"alias" structs:",underline"`
    ChangeMap map[string]interface{} `json:"-" structs:"-"`
}

func (e Example) TableName() string {
    return "example"
}
