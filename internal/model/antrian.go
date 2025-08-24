package model

import (
	"time"

	"gorm.io/gorm"
)

type Antrian struct {
	ID           int            `json:"id,omitempty" gorm:"primaryKey;column:id_antrian"`
	JadwalID     int            `json:"jadwal_id" gorm:"column:id_jadwal"`
	PasienID     int            `json:"pasien_id" gorm:"column:id_pasien"`
	NomorAntrian string         `json:"nomor_antrian" gorm:"column:nomor_antrian"`
	Prioritas    string         `json:"prioritas" gorm:"column:prioritas"`
	Status       string         `json:"status" gorm:"column:status"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index;column:deleted_at"`
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"column:updated_at"`

	Jadwal Jadwal `json:"jadwal" gorm:"foreignKey:JadwalID"`
	Pasien Pasien `json:"pasien" gorm:"foreignKey:PasienID"`
}

func (Antrian) TableName() string { return "antrian" }

type CreateAntrianRequest struct {
	JadwalID  int    `json:"jadwal_id" binding:"required,gt=0"`
	PasienID  int    `json:"pasien_id" binding:"required,gt=0"`
	Prioritas string `json:"prioritas" binding:"required,oneof=Gawat 'Non Gawat'"`
}

func (req *CreateAntrianRequest) ToModel(nomorAntrian string) Antrian {
	return Antrian{
		JadwalID:     req.JadwalID,
		PasienID:     req.PasienID,
		Prioritas:    req.Prioritas,
		Status:       "Menunggu",
		NomorAntrian: nomorAntrian,
	}
}

type UpdateAntrianRequest struct {
	Prioritas string `json:"prioritas" binding:"required,oneof=Gawat 'Non Gawat'"`
	Status    string `json:"status" binding:"required,oneof=Menunggu 'Menunggu Diagnosis' Selesai"`
}

func (req *UpdateAntrianRequest) ToModel() Antrian {
	return Antrian{
		Prioritas: req.Prioritas,
		Status:    req.Status,
	}
}

type AntrianResponse struct {
	ID           int    `json:"id"`
	NomorAntrian string `json:"nomor_antrian"`
	Prioritas    string `json:"prioritas"`
	Status       string `json:"status"`
	Jadwal       struct {
		ID      int    `json:"id"`
		Tanggal string `json:"tanggal"`
		Poli    struct {
			Nama string `json:"nama"`
		} `json:"poli"`
		Dokter struct {
			Nama string `json:"nama"`
		} `json:"dokter"`
	} `json:"jadwal"`
	Pasien struct {
		ID   int    `json:"id"`
		Nama string `json:"nama"`
	} `json:"pasien"`
}

func ToAntrianResponse(a Antrian) AntrianResponse {
	return AntrianResponse{
		ID:           a.ID,
		NomorAntrian: a.NomorAntrian,
		Prioritas:    a.Prioritas,
		Status:       a.Status,
		Jadwal: struct {
			ID      int    `json:"id"`
			Tanggal string `json:"tanggal"`
			Poli    struct {
				Nama string `json:"nama"`
			} `json:"poli"`
			Dokter struct {
				Nama string `json:"nama"`
			} `json:"dokter"`
		}{
			ID:      a.Jadwal.ID,
			Tanggal: a.Jadwal.Tanggal.Format("2006-01-02"),
			Poli: struct {
				Nama string `json:"nama"`
			}{Nama: a.Jadwal.Poli.Nama},
			Dokter: struct {
				Nama string `json:"nama"`
			}{Nama: a.Jadwal.Petugas.Nama},
		},
		Pasien: struct {
			ID   int    `json:"id"`
			Nama string `json:"nama"`
		}{
			ID:   a.Pasien.ID,
			Nama: a.Pasien.NamaPasien,
		},
	}
}

func ToAntrianResponseList(antrians []Antrian) []AntrianResponse {
	var responses []AntrianResponse
	for _, a := range antrians {
		responses = append(responses, ToAntrianResponse(a))
	}
	return responses
}
