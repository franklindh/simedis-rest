package model

import (
	"time"

	"gorm.io/gorm"
)

type Poli struct {
	ID        int            `json:"id,omitempty" gorm:"primaryKey;column:id_poli"`
	Nama      string         `json:"nama" gorm:"column:nama_poli"`
	Status    string         `json:"status" gorm:"column:status_poli"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;column:deleted_at"`
}

func (Poli) TableName() string {
	return "poli"
}

type CreatePoliRequest struct {
	Name   string `json:"name" binding:"required,min=3,max=50"`
	Status string `json:"status" binding:"required,oneof=aktif nonaktif"`
}

type UpdatePoliRequest struct {
	Name   string `json:"name" binding:"required,min=3,max=50"`
	Status string `json:"status" binding:"required,oneof=aktif nonaktif"`
}

type PoliResponse struct {
	ID        int       `json:"id"`
	Nama      string    `json:"nama"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToPoliResponse(p Poli) PoliResponse {
	return PoliResponse{
		ID:        p.ID,
		Nama:      p.Nama,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func ToPoliResponseList(polis []Poli) []PoliResponse {
	var responses []PoliResponse
	for _, p := range polis {
		responses = append(responses, ToPoliResponse(p))
	}
	return responses
}

func (req *CreatePoliRequest) ToModel() Poli {
	return Poli{
		Nama:   req.Name,
		Status: req.Status,
	}
}

func (req *UpdatePoliRequest) ToModel() Poli {
	return Poli{
		Nama:   req.Name,
		Status: req.Status,
	}
}
