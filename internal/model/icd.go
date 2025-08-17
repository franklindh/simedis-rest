package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Icd struct {
	ID           int            `json:"id,omitempty" gorm:"primaryKey;column:id_icd"`
	KodeIcd      string         `json:"kode_icd" gorm:"column:kode_icd;unique" binding:"required"`
	NamaPenyakit string         `json:"nama_penyakit" gorm:"column:nama_penyakit" binding:"required,min=3"`
	Deskripsi    sql.NullString `json:"deskripsi,omitempty" gorm:"column:deskripsi_penyakit"`
	Status       string         `json:"status" gorm:"column:status_icd" binding:"required,oneof=aktif nonaktif"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index;column:deleted_at"`
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"column:updated_at"`
}

func (Icd) TableName() string {
	return "icd"
}
