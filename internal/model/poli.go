package model

import (
	"time"

	"gorm.io/gorm"
)

type Poli struct {
	ID        int            `json:"id,omitempty" gorm:"primaryKey;column:id_poli"`
	Name      string         `json:"name" gorm:"column:nama_poli"`
	Status    string         `json:"status" gorm:"column:status_poli"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;column:deleted_at"`
}

func (Poli) TableName() string {
	return "poli"
}
