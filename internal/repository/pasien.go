package repository

import (
	"errors"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/pkg/utils/pagination"
	"gorm.io/gorm"
)

type ParamsGetAllPasien struct {
	NameFilter   string `form:"name" binding:"omitempty,sanitize"`
	NIKFilter    string `form:"nik" binding:"omitempty,numeric"`
	NoRekamMedis string `form:"no_rekam_medis" binding:"omitempty,sanitize"`
	SortBy       string `form:"sort" binding:"omitempty,sanitize"`
	Page         int    `form:"page" binding:"omitempty,gt=0"`
	PageSize     int    `form:"pageSize" binding:"omitempty,gt=0"`
}

type PasienRepository struct {
	DB *gorm.DB
}

func NewPasienRepository(db *gorm.DB) *PasienRepository {
	return &PasienRepository{DB: db}
}

func (r *PasienRepository) Create(pasien model.Pasien) (model.Pasien, error) {
	result := r.DB.Create(&pasien)
	return pasien, result.Error
}

func (r *PasienRepository) GetAll(params ParamsGetAllPasien) ([]model.Pasien, pagination.Metadata, error) {
	var pasien []model.Pasien
	var totalRecords int64

	db := r.DB.Model(&model.Pasien{})

	if params.NameFilter != "" {
		db = db.Where("nama_pasien ILIKE ?", "%"+params.NameFilter+"%")
	}
	if params.NIKFilter != "" {
		db = db.Where("nik ILIKE ?", "%"+params.NIKFilter+"%")
	}
	if params.NoRekamMedis != "" {
		db = db.Where("no_rekam_medis ILIKE ?", "%"+params.NoRekamMedis+"%")
	}

	if err := db.Count(&totalRecords).Error; err != nil {
		return nil, pagination.Metadata{}, err
	}

	metadata := pagination.CalculateMetadata(int(totalRecords), params.Page, params.PageSize)

	sortWhiteList := map[string]string{
		"nama_asc":  "nama_pasien ASC",
		"nama_desc": "nama_pasien DESC",
	}

	orderByClause := "nama_pasien ASC"
	if sort, ok := sortWhiteList[params.SortBy]; ok {
		orderByClause = sort
	}
	db = db.Order(orderByClause)

	db = db.Limit(metadata.PageSize).Offset((metadata.CurrentPage - 1) * metadata.PageSize)

	if err := db.Find(&pasien).Error; err != nil {
		return nil, pagination.Metadata{}, err
	}
	return pasien, metadata, nil
}

func (r *PasienRepository) GetById(id int) (model.Pasien, error) {
	var pasien model.Pasien
	result := r.DB.First(&pasien, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Pasien{}, ErrNotFound
		}
		return model.Pasien{}, result.Error
	}
	return pasien, nil
}

func (r *PasienRepository) Update(id int, pasien model.Pasien) (model.Pasien, error) {

	result := r.DB.Model(&model.Pasien{}).Where("id_pasien = ?", id).Updates(&pasien)
	if result.Error != nil {
		return model.Pasien{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.Pasien{}, ErrNotFound
	}
	return r.GetById(id)
}

func (r *PasienRepository) Delete(id int) error {
	result := r.DB.Delete(&model.Pasien{}, id)
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}

func (r *PasienRepository) GetLastID() (int, error) {
	var lastID int

	result := r.DB.Model(&model.Pasien{}).Select("id_pasien").Order("id_pasien DESC").Limit(1).Scan(&lastID)
	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, result.Error
	}
	return lastID, nil
}
