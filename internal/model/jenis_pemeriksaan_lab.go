package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type JenisPemeriksaanLab struct {
	ID              int            `json:"id,omitempty" gorm:"primaryKey;column:id_jenis_pemeriksaan"`
	NamaPemeriksaan string         `json:"nama_pemeriksaan" gorm:"column:nama_pemeriksaan;unique"`
	Satuan          sql.NullString `json:"satuan" gorm:"column:satuan"`
	NilaiRujukan    sql.NullString `json:"nilai_rujukan" gorm:"column:nilai_rujukan"`
	Kriteria        sql.NullString `json:"kriteria" gorm:"column:kriteria"`
	CreatedAt       time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index;column:deleted_at"`
}

func (JenisPemeriksaanLab) TableName() string {
	return "jenis_pemeriksaan_lab"
}

type CreateJenisPemeriksaanLabRequest struct {
	NamaPemeriksaan string `json:"nama_pemeriksaan" binding:"required,min=3"`
	Satuan          string `json:"satuan,omitempty"`
	NilaiRujukan    string `json:"nilai_rujukan,omitempty"`
	Kriteria        string `json:"kriteria,omitempty"`
}

type UpdateJenisPemeriksaanLabRequest struct {
	NamaPemeriksaan string `json:"nama_pemeriksaan" binding:"required,min=3"`
	Satuan          string `json:"satuan,omitempty"`
	NilaiRujukan    string `json:"nilai_rujukan,omitempty"`
	Kriteria        string `json:"kriteria,omitempty"`
}
