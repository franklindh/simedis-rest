package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Petugas struct {
	ID        int            `json:"id,omitempty" gorm:"primaryKey;column:id_petugas"`
	PoliID    sql.NullInt64  `json:"poli_id" gorm:"column:id_poli"`
	Username  string         `json:"username" gorm:"column:username_petugas;unique"`
	Nama      string         `json:"nama" gorm:"column:nama_petugas"`
	Status    string         `json:"status" gorm:"column:status"`
	Role      string         `json:"role" gorm:"column:role"`
	Password  string         `json:"-" gorm:"column:password"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;column:deleted_at"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`
}

type PetugasResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Nama      string    `json:"nama"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	PoliID    *int64    `json:"poli_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePetugasRequest struct {
	Username string `json:"username" binding:"required,min=5,max=20,alphanum,sanitize"`
	Nama     string `json:"nama" binding:"required,min=3,max=50,sanitize"`
	Status   string `json:"status" binding:"required,oneof=aktif nonaktif"`
	Role     string `json:"role" binding:"required,oneof=Administrasi Poliklinik Dokter Lab"`
	PoliID   *int64 `json:"poli_id,omitempty"`
}

type UpdatePetugasRequest struct {
	Nama   string `json:"nama" binding:"required,min=3,max=50,sanitize"`
	Status string `json:"status" binding:"required,oneof=aktif nonaktif"`
	Role   string `json:"role" binding:"required,oneof=Administrasi Poliklinik Dokter Lab"`
	PoliID *int64 `json:"poli_id,omitempty"`
}

type LoginPetugasRequest struct {
	Username string `json:"username" binding:"required,sanitize"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

func (req *CreatePetugasRequest) ToModel() Petugas {
	petugas := Petugas{
		Username: req.Username,
		Nama:     req.Nama,
		Status:   req.Status,
		Role:     req.Role,
	}
	if req.PoliID != nil {
		petugas.PoliID = sql.NullInt64{Int64: *req.PoliID, Valid: true}
	}
	return petugas
}

func (req *UpdatePetugasRequest) ToModel() Petugas {
	petugas := Petugas{
		Nama:   req.Nama,
		Status: req.Status,
		Role:   req.Role,
	}
	if req.PoliID != nil {
		petugas.PoliID = sql.NullInt64{Int64: *req.PoliID, Valid: true}
	} else {

		petugas.PoliID = sql.NullInt64{Valid: false}
	}
	return petugas
}

func ToPetugasResponse(p Petugas) PetugasResponse {
	var poliID *int64
	if p.PoliID.Valid {
		poliID = &p.PoliID.Int64
	}

	return PetugasResponse{
		ID:        p.ID,
		Username:  p.Username,
		Nama:      p.Nama,
		Role:      p.Role,
		Status:    p.Status,
		PoliID:    poliID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func ToPetugasResponseList(petugasList []Petugas) []PetugasResponse {
	var responses []PetugasResponse
	for _, p := range petugasList {
		responses = append(responses, ToPetugasResponse(p))
	}
	return responses
}

func (Petugas) TableName() string {
	return "petugas"
}
