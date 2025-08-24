package repository

import (
	"database/sql"
	"errors"

	"github.com/franklindh/simedis-api/internal/model"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("data not found")

type PoliRepository struct {
	DB *gorm.DB
}

func NewPoliRepository(db *gorm.DB) *PoliRepository {
	return &PoliRepository{DB: db}
}

func (r *PoliRepository) GetAll() ([]model.Poli, error) {

	var poli []model.Poli
	result := r.DB.Find(&poli)
	return poli, result.Error
}

func (r *PoliRepository) GetById(id int) (model.Poli, error) {

	var poli model.Poli
	result := r.DB.First(&poli, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Poli{}, ErrNotFound
		}
		return model.Poli{}, result.Error
	}
	return poli, nil
}

func (r *PoliRepository) Create(poli model.Poli) (model.Poli, error) {

	result := r.DB.Create(&poli)
	return poli, result.Error
}

func (r *PoliRepository) Update(id int, poli model.Poli) (model.Poli, error) {
	result := r.DB.Model(&model.Poli{}).Where("id_poli = ?", id).Updates(&poli)
	if result.Error != nil {
		return model.Poli{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.Poli{}, ErrNotFound
	}

	return r.GetById(id)
}

func (r *PoliRepository) Delete(id int) error {

	result := r.DB.Delete(&model.Poli{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PoliRepository) FindByNameIncludingDeleted(nama string) (model.Poli, error) {

	var poli model.Poli

	result := r.DB.Unscoped().Where("nama_poli = ?", nama).First(&poli)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Poli{}, sql.ErrNoRows
		}
		return model.Poli{}, result.Error
	}
	return poli, nil
}

func (r *PoliRepository) FindByName(name string) (model.Poli, error) {
	var poli model.Poli
	result := r.DB.Where("nama_poli = ?", name).First(&poli)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Poli{}, ErrNotFound
		}
		return model.Poli{}, result.Error
	}
	return poli, nil
}
