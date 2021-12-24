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
    Name      string `json:"name"`
    Alias     string `json:"alias"`
    changeMap map[string]interface{}
}

func (e Example) TableName() string {
    return "example"
}

func (e *Example) GetChangeMap() map[string]interface{} {
    return e.changeMap
}
