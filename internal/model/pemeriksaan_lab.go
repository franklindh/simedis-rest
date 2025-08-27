package model

import "time"

type PemeriksaanLab struct {
	ID                  int                 `json:"id,omitempty" gorm:"primaryKey;column:id_pemeriksaan_lab"`
	PemeriksaanID       int                 `json:"pemeriksaan_id" gorm:"column:id_pemeriksaan"`
	JenisPemeriksaanID  int                 `json:"jenis_pemeriksaan_id" gorm:"column:id_jenis_pemeriksaan"`
	Hasil               string              `json:"hasil" gorm:"column:hasil"`
	CreatedAt           time.Time           `json:"created_at" gorm:"column:created_at"`
	UpdatedAt           time.Time           `json:"updated_at" gorm:"column:updated_at"`
	JenisPemeriksaanLab JenisPemeriksaanLab `json:"jenis_pemeriksaan" gorm:"foreignKey:JenisPemeriksaanID"`
}

func (PemeriksaanLab) TableName() string {
	return "pemeriksaan_lab"
}

type CreateHasilLabRequest struct {
	JenisPemeriksaanID int    `json:"jenis_pemeriksaan_id" binding:"required,gt=0"`
	Hasil              string `json:"hasil" binding:"required,sanitize"`
}

func (req *CreateHasilLabRequest) ToModel(pemeriksaanID int) PemeriksaanLab {
	return PemeriksaanLab{
		PemeriksaanID:      pemeriksaanID,
		JenisPemeriksaanID: req.JenisPemeriksaanID,
		Hasil:              req.Hasil,
	}
}

type UpdateHasilLabRequest struct {
	Hasil string `json:"hasil" binding:"required,sanitize"`
}

func (req *UpdateHasilLabRequest) ToModel() PemeriksaanLab {
	return PemeriksaanLab{
		Hasil: req.Hasil,
	}
}

type PemeriksaanLabResponse struct {
	ID               int    `json:"id"`
	Hasil            string `json:"hasil"`
	JenisPemeriksaan struct {
		ID   int    `json:"id"`
		Nama string `json:"nama"`
	} `json:"jenis_pemeriksaan"`
}

func ToPemeriksaanLabResponse(p PemeriksaanLab) PemeriksaanLabResponse {
	return PemeriksaanLabResponse{
		ID:    p.ID,
		Hasil: p.Hasil,
		JenisPemeriksaan: struct {
			ID   int    `json:"id"`
			Nama string `json:"nama"`
		}{
			ID:   p.JenisPemeriksaanLab.ID,
			Nama: p.JenisPemeriksaanLab.NamaPemeriksaan,
		},
	}
}

func ToPemeriksaanLabResponseList(pemeriksaanLabs []PemeriksaanLab) []PemeriksaanLabResponse {
	var responses []PemeriksaanLabResponse
	for _, p := range pemeriksaanLabs {
		responses = append(responses, ToPemeriksaanLabResponse(p))
	}
	return responses
}
