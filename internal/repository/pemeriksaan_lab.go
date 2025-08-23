package repository

import (
	"errors"

	"github.com/franklindh/simedis-api/internal/model"
	"gorm.io/gorm"
)

type PemeriksaanLabRepository struct {
	DB *gorm.DB
}

func NewPemeriksaanLabRepository(db *gorm.DB) *PemeriksaanLabRepository {
	return &PemeriksaanLabRepository{DB: db}
}

func (r *PemeriksaanLabRepository) Create(hasilLab model.PemeriksaanLab) (model.PemeriksaanLab, error) {
	result := r.DB.Create(&hasilLab)
	if result.Error != nil {
		return model.PemeriksaanLab{}, result.Error
	}
	return r.GetById(hasilLab.ID)
}

func (r *PemeriksaanLabRepository) GetAllByPemeriksaanID(pemeriksaanID int) ([]model.PemeriksaanLab, error) {
	var results []model.PemeriksaanLab
	err := r.DB.Preload("JenisPemeriksaanLab").Where("id_pemeriksaan = ?", pemeriksaanID).Find(&results).Error
	return results, err
}

func (r *PemeriksaanLabRepository) GetById(id int) (model.PemeriksaanLab, error) {
	var hasilLab model.PemeriksaanLab
	result := r.DB.Preload("JenisPemeriksaanLab").First(&hasilLab, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.PemeriksaanLab{}, ErrNotFound
		}
		return model.PemeriksaanLab{}, result.Error
	}
	return hasilLab, nil
}

func (r *PemeriksaanLabRepository) Update(id int, hasilLab model.PemeriksaanLab) (model.PemeriksaanLab, error) {
	result := r.DB.Model(&model.PemeriksaanLab{}).Where("id_pemeriksaan_lab = ?", id).Updates(&hasilLab)
	if result.Error != nil {
		return model.PemeriksaanLab{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.PemeriksaanLab{}, ErrNotFound
	}
	return r.GetById(id)
}

func (r *PemeriksaanLabRepository) Delete(id int) error {
	result := r.DB.Delete(&model.PemeriksaanLab{}, id)
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}
