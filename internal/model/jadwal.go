package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Jadwal struct {
	ID           int            `json:"id,omitempty" gorm:"primaryKey;column:id_jadwal"`
	PetugasID    int            `json:"petugas_id" gorm:"column:id_petugas"`
	PoliID       int            `json:"poli_id" gorm:"column:id_poli"`
	Tanggal      string         `json:"tanggal" gorm:"column:tanggal_praktik"`
	WaktuMulai   string         `json:"waktu_mulai" gorm:"column:waktu_mulai"`
	WaktuSelesai string         `json:"waktu_selesai" gorm:"column:waktu_selesai"`
	Keterangan   sql.NullString `json:"keterangan" gorm:"column:keterangan"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index;column:deleted_at"`
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"column:updated_at"`

	Petugas Petugas `json:"petugas" gorm:"foreignKey:PetugasID"`
	Poli    Poli    `json:"poli" gorm:"foreignKey:PoliID"`
}

func (Jadwal) TableName() string {
	return "jadwal"
}

type JadwalRequest struct {
	PetugasID    int    `json:"petugas_id" binding:"required,gt=0"`
	PoliID       int    `json:"poli_id" binding:"required,gt=0"`
	Tanggal      string `json:"tanggal" binding:"required,datetime=2006-01-02"`
	WaktuMulai   string `json:"waktu_mulai" binding:"required,datetime=15:04"`
	WaktuSelesai string `json:"waktu_selesai" binding:"required,datetime=15:04"`
	Keterangan   string `json:"keterangan,omitempty"`
}

type JadwalDetail struct {
	ID           int            `json:"id"`
	Tanggal      string         `json:"tanggal"`
	WaktuMulai   string         `json:"waktu_mulai"`
	WaktuSelesai string         `json:"waktu_selesai"`
	Keterangan   sql.NullString `json:"keterangan"`
	Petugas      struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"petugas"`
	Poli struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"poli"`
}

func (req *JadwalRequest) ToModel() Jadwal {
	return Jadwal{
		PetugasID:    req.PetugasID,
		PoliID:       req.PoliID,
		Tanggal:      req.Tanggal,
		WaktuMulai:   req.WaktuMulai,
		WaktuSelesai: req.WaktuSelesai,
		Keterangan:   sql.NullString{String: req.Keterangan, Valid: req.Keterangan != ""},
	}
}
