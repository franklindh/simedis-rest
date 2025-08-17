package repository

import (
	"errors"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/pkg/utils/pagination"
	"gorm.io/gorm"
)

type ParamsGetAllIcd struct {
	KodeFilter   string
	NamaFilter   string
	StatusFilter string
	SortBy       string
	Page         int
	PageSize     int
}

type IcdRepository struct {
	DB *gorm.DB
}

func NewIcdRepository(db *gorm.DB) *IcdRepository {
	return &IcdRepository{DB: db}
}

func (r *IcdRepository) Create(icd model.Icd) (model.Icd, error) {
	result := r.DB.Create(&icd)
	return icd, result.Error
}

func (r *IcdRepository) GetAll(params ParamsGetAllIcd) ([]model.Icd, pagination.Metadata, error) {
	var allIcd []model.Icd
	var totalRecords int64

	db := r.DB.Model(&model.Icd{})

	if params.KodeFilter != "" {
		db = db.Where("kode_icd ILIKE ?", "%"+params.KodeFilter+"%")
	}
	if params.NamaFilter != "" {
		db = db.Where("nama_penyakit ILIKE ?", "%"+params.NamaFilter+"%")
	}
	if params.StatusFilter != "" {
		db = db.Where("status_icd = ?", params.StatusFilter)
	}

	if err := db.Count(&totalRecords).Error; err != nil {
		return nil, pagination.Metadata{}, err
	}

	metadata := pagination.CalculateMetadata(int(totalRecords), params.Page, params.PageSize)

	sortWhitelist := map[string]string{
		"kode_asc":  "kode_icd ASC",
		"kode_desc": "kode_icd DESC",
		"nama_asc":  "nama_penyakit ASC",
		"nama_desc": "nama_penyakit DESC",
	}
	orderByClause := "kode_icd ASC"
	if sort, ok := sortWhitelist[params.SortBy]; ok {
		orderByClause = sort
	}
	db = db.Order(orderByClause)

	db = db.Limit(metadata.PageSize).Offset((metadata.CurrentPage - 1) * metadata.PageSize)

	if err := db.Find(&allIcd).Error; err != nil {
		return nil, pagination.Metadata{}, err
	}

	return allIcd, metadata, nil
}

func (r *IcdRepository) GetByID(id int) (model.Icd, error) {
	var icd model.Icd
	result := r.DB.First(&icd, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Icd{}, ErrNotFound
		}
		return model.Icd{}, result.Error
	}
	return icd, nil
}

func (r *IcdRepository) Update(id int, icd model.Icd) (model.Icd, error) {
	icd.ID = id
	result := r.DB.Model(&model.Icd{}).Where("id_icd = ?", id).Updates(&icd)
	if result.Error != nil {
		return model.Icd{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.Icd{}, ErrNotFound
	}
	return r.GetByID(id)
}

func (r *IcdRepository) Delete(id int) error {
	result := r.DB.Delete(&model.Icd{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
