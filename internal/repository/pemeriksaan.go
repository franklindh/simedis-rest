package repository

import (
	"errors"

	"github.com/franklindh/simedis-api/internal/model"
	"gorm.io/gorm"
)

type PemeriksaanRepository struct {
	DB *gorm.DB
}

func NewPemeriksaanRepository(db *gorm.DB) *PemeriksaanRepository {
	return &PemeriksaanRepository{DB: db}
}

func (r *PemeriksaanRepository) CheckExistingPemeriksaan(antrianID int) error {
	var count int64
	result := r.DB.Model(&model.Pemeriksaan{}).Where("id_antrian = ?", antrianID).Count(&count)
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return nil
	}
	return ErrNotFound
}

func (r *PemeriksaanRepository) Create(pemeriksaan model.Pemeriksaan) (model.Pemeriksaan, error) {
	result := r.DB.Create(&pemeriksaan)
	if result.Error != nil {
		return model.Pemeriksaan{}, result.Error
	}

	return r.GetById(pemeriksaan.ID)
}

func (r *PemeriksaanRepository) GetById(id int) (model.Pemeriksaan, error) {
	var pemeriksaan model.Pemeriksaan
	result := r.DB.
		Preload("Icd").
		Preload("Antrian.Pasien").
		Preload("Antrian.Jadwal.Petugas").
		Preload("Antrian.Jadwal.Poli").
		First(&pemeriksaan, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Pemeriksaan{}, ErrNotFound
		}
		return model.Pemeriksaan{}, result.Error
	}
	return pemeriksaan, nil
}

func (r *PemeriksaanRepository) GetAllByPasienID(pasienID int) ([]model.Pemeriksaan, error) {
	var allPemeriksaan []model.Pemeriksaan
	result := r.DB.
		Joins("JOIN antrian ON pemeriksaan.id_antrian = antrian.id_antrian").
		Where("antrian.id_pasien = ?", pasienID).
		Preload("Icd").
		Preload("Antrian.Pasien").
		Preload("Antrian.Jadwal.Petugas").
		Preload("Antrian.Jadwal.Poli").
		Order("tanggal_pemeriksaan DESC, created_at DESC").
		Find(&allPemeriksaan)

	return allPemeriksaan, result.Error
}

func (r *PemeriksaanRepository) Update(id int, pemeriksaan model.Pemeriksaan) (model.Pemeriksaan, error) {
	result := r.DB.Model(&model.Pemeriksaan{}).Where("id_pemeriksaan = ?", id).Updates(&pemeriksaan)
	if result.Error != nil {
		return model.Pemeriksaan{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.Pemeriksaan{}, ErrNotFound
	}
	return r.GetById(id)
}

func (r *PemeriksaanRepository) Delete(id int) error {
	result := r.DB.Delete(&model.Pemeriksaan{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
