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
	Tanggal      time.Time      `json:"tanggal" gorm:"column:tanggal_praktik"`
	WaktuMulai   time.Time      `json:"waktu_mulai" gorm:"column:waktu_mulai"`
	WaktuSelesai time.Time      `json:"waktu_selesai" gorm:"column:waktu_selesai"`
	Keterangan   sql.NullString `json:"keterangan" gorm:"column:keterangan"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index;column:deleted_at"`
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"column:updated_at"`
	Petugas      Petugas        `json:"petugas" gorm:"foreignKey:PetugasID"`
	Poli         Poli           `json:"poli" gorm:"foreignKey:PoliID"`
}

func (Jadwal) TableName() string { return "jadwal" }

type JadwalRequest struct {
	PetugasID    int    `json:"petugas_id" binding:"required,gt=0"`
	PoliID       int    `json:"poli_id" binding:"required,gt=0"`
	Tanggal      string `json:"tanggal" binding:"required,datetime=2006-01-02"`
	WaktuMulai   string `json:"waktu_mulai" binding:"required,datetime=15:04"`
	WaktuSelesai string `json:"waktu_selesai" binding:"required,datetime=15:04"`
	Keterangan   string `json:"keterangan,omitempty"`
}

func (req *JadwalRequest) ToModel() Jadwal {
	parsedTanggal, _ := time.Parse("2006-01-02", req.Tanggal)
	parsedMulai, _ := time.Parse("15:04", req.WaktuMulai)
	parsedSelesai, _ := time.Parse("15:04", req.WaktuSelesai)
	return Jadwal{
		PetugasID:    req.PetugasID,
		PoliID:       req.PoliID,
		Tanggal:      parsedTanggal,
		WaktuMulai:   parsedMulai,
		WaktuSelesai: parsedSelesai,
		Keterangan:   sql.NullString{String: req.Keterangan, Valid: req.Keterangan != ""},
	}
}

type JadwalResponse struct {
	ID           int    `json:"id"`
	Tanggal      string `json:"tanggal"`
	WaktuMulai   string `json:"waktu_mulai"`
	WaktuSelesai string `json:"waktu_selesai"`
	Keterangan   string `json:"keterangan,omitempty"`
	Petugas      struct {
		ID   int    `json:"id"`
		Nama string `json:"nama"`
	} `json:"petugas"`
	Poli struct {
		ID   int    `json:"id"`
		Nama string `json:"nama"`
	} `json:"poli"`
}

func ToJadwalResponse(j Jadwal) JadwalResponse {
	return JadwalResponse{
		ID:           j.ID,
		Tanggal:      j.Tanggal.Format("2006-01-02"),
		WaktuMulai:   j.WaktuMulai.Format("15:04"),
		WaktuSelesai: j.WaktuSelesai.Format("15:04"),
		Keterangan:   j.Keterangan.String,
		Petugas: struct {
			ID   int    `json:"id"`
			Nama string `json:"nama"`
		}{
			ID:   j.Petugas.ID,
			Nama: j.Petugas.Nama,
		},
		Poli: struct {
			ID   int    `json:"id"`
			Nama string `json:"nama"`
		}{
			ID:   j.Poli.ID,
			Nama: j.Poli.Nama,
		},
	}
}

func ToJadwalResponseList(jadwals []Jadwal) []JadwalResponse {
	var responses []JadwalResponse
	for _, j := range jadwals {
		responses = append(responses, ToJadwalResponse(j))
	}
	return responses
}
