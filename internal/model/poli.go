package model

import (
	"time"

	"gorm.io/gorm"
)

type Poli struct {
	ID        int            `json:"id,omitempty" gorm:"primaryKey;column:id_poli"`
	Nama      string         `json:"nama" gorm:"column:nama_poli" binding:"required,min=3,max=50"`
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
