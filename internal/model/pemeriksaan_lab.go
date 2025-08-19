package model

import "time"

type PemeriksaanLab struct {
	ID                 int       `json:"id,omitempty" gorm:"primaryKey;column:id_pemeriksaan_lab"`
	PemeriksaanID      int       `json:"pemeriksaan_id" gorm:"column:id_pemeriksaan"`
	JenisPemeriksaanID int       `json:"jenis_pemeriksaan_id" gorm:"column:id_jenis_pemeriksaan"`
	Hasil              string    `json:"hasil" gorm:"column:hasil"`
	CreatedAt          time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"column:updated_at"`

	JenisPemeriksaanLab JenisPemeriksaanLab `json:"jenis_pemeriksaan,omitempty" gorm:"foreignKey:JenisPemeriksaanID"`
}

func (PemeriksaanLab) TableName() string {
	return "pemeriksaan_lab"
}

type CreateHasilLabRequest struct {
	JenisPemeriksaanID int    `json:"jenis_pemeriksaan_id" binding:"required,gt=0"`
	Hasil              string `json:"hasil" binding:"required"`
}

type UpdateHasilLabRequest struct {
	Hasil string `json:"hasil" binding:"required"`
}
