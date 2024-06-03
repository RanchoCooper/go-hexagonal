package dto

import (
	"time"
)

type CreateExampleReq struct {
	Name  string `json:"name" binding:"required" message:"name is a required field"`
	Alias string `json:"alias" binding:"required" message:"alias is a required field"`
}

type CreateExampleResp struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	Alias     string    `json:"alias"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeleteExampleReq struct {
	Id int `uri:"id" binding:"required"`
}

type UpdateExampleReq struct {
	Id    uint   `uri:"id"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

type GetExampleReq struct {
	Id int `uri:"id" binding:"required"`
}

type GetExampleResponse struct {
	Id        int       `uri:"id"`
	Name      string    `json:"name"`
	Alias     string    `json:"alias"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
