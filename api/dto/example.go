package dto

import (
	"time"
)

type CreateExampleReq struct {
	Name  string `json:"name" binding:"required,min=1,max=100" message:"name is required and must be 1-100 characters"`
	Alias string `json:"alias" binding:"required,min=1,max=100" message:"alias is required and must be 1-100 characters"`
}

type CreateExampleResp struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	Alias     string    `json:"alias"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeleteExampleReq struct {
	Id int `uri:"id" binding:"required,min=1"`
}

type UpdateExampleReq struct {
	Id    uint   `uri:"id" binding:"required"`
	Name  string `json:"name" binding:"omitempty,min=1,max=100"`
	Alias string `json:"alias" binding:"omitempty,min=1,max=100"`
}

type GetExampleReq struct {
	Id int `uri:"id" binding:"required,min=1"`
}

type GetExampleResponse struct {
	Id        int       `uri:"id"`
	Name      string    `json:"name"`
	Alias     string    `json:"alias"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
