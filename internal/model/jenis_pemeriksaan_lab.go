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

func (req *CreateJenisPemeriksaanLabRequest) ToModel() JenisPemeriksaanLab {
	return JenisPemeriksaanLab{
		NamaPemeriksaan: req.NamaPemeriksaan,
		Satuan:          sql.NullString{String: req.Satuan, Valid: req.Satuan != ""},
		NilaiRujukan:    sql.NullString{String: req.NilaiRujukan, Valid: req.NilaiRujukan != ""},
		Kriteria:        sql.NullString{String: req.Kriteria, Valid: req.Kriteria != ""},
	}
}

type UpdateJenisPemeriksaanLabRequest struct {
	NamaPemeriksaan string `json:"nama_pemeriksaan" binding:"required,min=3"`
	Satuan          string `json:"satuan,omitempty"`
	NilaiRujukan    string `json:"nilai_rujukan,omitempty"`
	Kriteria        string `json:"kriteria,omitempty"`
}

func (req *UpdateJenisPemeriksaanLabRequest) ToModel() JenisPemeriksaanLab {
	return JenisPemeriksaanLab{
		NamaPemeriksaan: req.NamaPemeriksaan,
		Satuan:          sql.NullString{String: req.Satuan, Valid: req.Satuan != ""},
		NilaiRujukan:    sql.NullString{String: req.NilaiRujukan, Valid: req.NilaiRujukan != ""},
		Kriteria:        sql.NullString{String: req.Kriteria, Valid: req.Kriteria != ""},
	}
}

type JenisPemeriksaanLabResponse struct {
	ID              int    `json:"id"`
	NamaPemeriksaan string `json:"nama_pemeriksaan"`
	Satuan          string `json:"satuan,omitempty"`
	NilaiRujukan    string `json:"nilai_rujukan,omitempty"`
	Kriteria        string `json:"kriteria,omitempty"`
}

func ToJenisPemeriksaanLabResponse(jpl JenisPemeriksaanLab) JenisPemeriksaanLabResponse {
	return JenisPemeriksaanLabResponse{
		ID:              jpl.ID,
		NamaPemeriksaan: jpl.NamaPemeriksaan,
		Satuan:          jpl.Satuan.String,
		NilaiRujukan:    jpl.NilaiRujukan.String,
		Kriteria:        jpl.Kriteria.String,
	}
}

func ToJenisPemeriksaanLabResponseList(list []JenisPemeriksaanLab) []JenisPemeriksaanLabResponse {
	var responses []JenisPemeriksaanLabResponse
	for _, jpl := range list {
		responses = append(responses, ToJenisPemeriksaanLabResponse(jpl))
	}
	return responses
}
