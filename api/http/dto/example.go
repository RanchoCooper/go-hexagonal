package dto

import (
    "time"
)

/**
 * @author Rancho
 * @date 2022/1/6
 */

type CreateExampleReq struct {
    Name  string `json:"name" validate:"required"`
    Alias string `json:"alias" validate:"required"`
}

type CreateExampleResp struct {
    Id        uint      `json:"id"`
    Name      string    `json:"name"`
    Alias     string    `json:"alias"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type DeleteExampleReq struct {
    Id int `uri:"id" validate:"required"`
}

type UpdateExampleReq struct {
    Id    uint   `uri:"id"`
    Name  string `json:"name"`
    Alias string `json:"alias"`
}
