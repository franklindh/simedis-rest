package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Petugas struct {
	ID                 int            `json:"id,omitempty" gorm:"primaryKey;column:id_petugas"`
	PoliID             sql.NullInt64  `json:"poli_id" gorm:"column:id_poli"`
	Username           string         `json:"username" gorm:"column:username_petugas;unique"`
	Name               string         `json:"name" gorm:"column:nama_petugas"`
	Status             string         `json:"status" gorm:"column:status"`
	Role               string         `json:"role" gorm:"column:role"`
	Password           string         `json:"-" gorm:"column:password"`
	WajibGantiPassword bool           `json:"wajib_ganti_password,omitempty" gorm:"column:wajib_ganti_password"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index;column:deleted_at"`
	CreatedAt          time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt          time.Time      `json:"updated_at" gorm:"column:updated_at"`
}

func (Petugas) TableName() string {
	return "petugas"
}
