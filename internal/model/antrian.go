package model

import (
	"time"

	"gorm.io/gorm"
)

type Antrian struct {
	ID           int            `json:"id,omitempty" gorm:"primaryKey;column:id_antrian"`
	JadwalID     int            `json:"jadwal_id" gorm:"column:id_jadwal" binding:"required,gt=0"`
	PasienID     int            `json:"pasien_id" gorm:"column:id_pasien" binding:"required,gt=0"`
	NomorAntrian string         `json:"nomor_antrian" gorm:"column:nomor_antrian"`
	Prioritas    string         `json:"prioritas" gorm:"column:prioritas" binding:"required,oneof=Gawat 'Non Gawat'"`
	Status       string         `json:"status" gorm:"column:status" binding:"required,oneof=Menunggu 'Menunggu Diagnosis' Selesai"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index;column:deleted_at"`
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"column:updated_at"`

	Jadwal Jadwal `json:"jadwal,omitempty" gorm:"foreignKey:JadwalID"`
	Pasien Pasien `json:"pasien,omitempty" gorm:"foreignKey:PasienID"`
}

type AntrianCreateInput struct {
	JadwalID  int    `json:"jadwal_id" binding:"required,gt=0"`
	PasienID  int    `json:"pasien_id" binding:"required,gt=0"`
	Prioritas string `json:"prioritas" binding:"required,oneof=Gawat 'Non Gawat'"`
	Status    string `json:"status" binding:"required,oneof=Menunggu 'Menunggu Diagnosis' Selesai"`
}

type AntrianDetail struct {
	ID           int    `json:"id"`
	NomorAntrian string `json:"nomor_antrian"`
	Prioritas    string `json:"prioritas"`
	Status       string `json:"status"`
	Jadwal       struct {
		ID      int    `json:"id"`
		Tanggal string `json:"tanggal"`
		Poli    struct {
			Name string `json:"name"`
		} `json:"poli"`
		Dokter struct {
			Name string `json:"name"`
		} `json:"dokter"`
	} `json:"jadwal"`
	Pasien struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"pasien"`
}

func (Antrian) TableName() string {
	return "antrian"
}
