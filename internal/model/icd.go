package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Icd struct {
	ID           int            `json:"id,omitempty" gorm:"primaryKey;column:id_icd"`
	KodeIcd      string         `json:"kode_icd" gorm:"column:kode_icd;unique"`
	NamaPenyakit string         `json:"nama_penyakit" gorm:"column:nama_penyakit"`
	Deskripsi    sql.NullString `json:"deskripsi,omitempty" gorm:"column:deskripsi_penyakit"`
	Status       string         `json:"status" gorm:"column:status_icd"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index;column:deleted_at"`
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"column:updated_at"`
}

func (Icd) TableName() string {
	return "icd"
}

type CreateIcdRequest struct {
	KodeIcd      string `json:"kode_icd" binding:"required"`
	NamaPenyakit string `json:"nama_penyakit" binding:"required,min=3"`
	Deskripsi    string `json:"deskripsi,omitempty"`
	Status       string `json:"status" binding:"required,oneof=aktif nonaktif"`
}

func (req *CreateIcdRequest) ToModel() Icd {
	return Icd{
		KodeIcd:      req.KodeIcd,
		NamaPenyakit: req.NamaPenyakit,
		Deskripsi:    sql.NullString{String: req.Deskripsi, Valid: req.Deskripsi != ""},
		Status:       req.Status,
	}
}

type UpdateIcdRequest struct {
	KodeIcd      string `json:"kode_icd" binding:"required"`
	NamaPenyakit string `json:"nama_penyakit" binding:"required,min=3"`
	Deskripsi    string `json:"deskripsi,omitempty"`
	Status       string `json:"status" binding:"required,oneof=aktif nonaktif"`
}

func (req *UpdateIcdRequest) ToModel() Icd {
	return Icd{
		KodeIcd:      req.KodeIcd,
		NamaPenyakit: req.NamaPenyakit,
		Deskripsi:    sql.NullString{String: req.Deskripsi, Valid: req.Deskripsi != ""},
		Status:       req.Status,
	}
}

type IcdResponse struct {
	ID           int    `json:"id"`
	KodeIcd      string `json:"kode_icd"`
	NamaPenyakit string `json:"nama_penyakit"`
	Deskripsi    string `json:"deskripsi,omitempty"`
	Status       string `json:"status"`
}

func ToIcdResponse(icd Icd) IcdResponse {
	return IcdResponse{
		ID:           icd.ID,
		KodeIcd:      icd.KodeIcd,
		NamaPenyakit: icd.NamaPenyakit,
		Deskripsi:    icd.Deskripsi.String,
		Status:       icd.Status,
	}
}

func ToIcdResponseList(icds []Icd) []IcdResponse {
	var responses []IcdResponse
	for _, icd := range icds {
		responses = append(responses, ToIcdResponse(icd))
	}
	return responses
}
