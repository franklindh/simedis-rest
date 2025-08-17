package model

import (
	"time"

	"gorm.io/gorm"
)

type Poli struct {
	ID        int            `json:"id,omitempty" gorm:"primaryKey;column:id_poli"`
	Nama      string         `json:"nama" gorm:"column:nama_poli" binding:"required,min=3,max=50"`
	Status    string         `json:"status" gorm:"column:status_poli" binding:"required,oneof=aktif nonaktif"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;column:deleted_at"`
}

func (Poli) TableName() string {
	return "poli"
}
