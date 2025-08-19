package repository

import (
	"errors"

	"github.com/franklindh/simedis-api/internal/model"
	"gorm.io/gorm"
)

type JenisPemeriksaanLabRepository struct {
	DB *gorm.DB
}

func NewJenisPemeriksaanLabRepository(db *gorm.DB) *JenisPemeriksaanLabRepository {
	return &JenisPemeriksaanLabRepository{DB: db}
}

func (r *JenisPemeriksaanLabRepository) Create(jenis model.JenisPemeriksaanLab) (model.JenisPemeriksaanLab, error) {
	result := r.DB.Create(&jenis)
	return jenis, result.Error
}

func (r *JenisPemeriksaanLabRepository) GetAll() ([]model.JenisPemeriksaanLab, error) {
	var allJenis []model.JenisPemeriksaanLab
	result := r.DB.Find(&allJenis)
	return allJenis, result.Error
}

func (r *JenisPemeriksaanLabRepository) GetByID(id int) (model.JenisPemeriksaanLab, error) {
	var jenis model.JenisPemeriksaanLab
	result := r.DB.First(&jenis, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.JenisPemeriksaanLab{}, ErrNotFound
		}
		return model.JenisPemeriksaanLab{}, result.Error
	}
	return jenis, nil
}

func (r *JenisPemeriksaanLabRepository) Update(id int, jenis model.JenisPemeriksaanLab) (model.JenisPemeriksaanLab, error) {
	result := r.DB.Model(&model.JenisPemeriksaanLab{}).Where("id_jenis_pemeriksaan = ?", id).Updates(&jenis)
	if result.Error != nil {
		return model.JenisPemeriksaanLab{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.JenisPemeriksaanLab{}, ErrNotFound
	}
	return r.GetByID(id)
}

func (r *JenisPemeriksaanLabRepository) Delete(id int) error {
	result := r.DB.Delete(&model.JenisPemeriksaanLab{}, id)
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}
