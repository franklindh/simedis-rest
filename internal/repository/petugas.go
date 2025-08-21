package repository

import (
	"errors"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/pkg/utils/pagination"
	"gorm.io/gorm"
)

type PetugasRepository interface {
	Create(petugas model.Petugas) (model.Petugas, error)
	GetAll(params ParamsGetAllPetugas) ([]model.Petugas, pagination.Metadata, error)
	GetByID(id int) (model.Petugas, error)
	GetByUsername(username string) (model.Petugas, error)
	Update(id int, petugas model.Petugas) (model.Petugas, error)
	Delete(id int) error
}

type ParamsGetAllPetugas struct {
	NameOrUsernameFilter string
	RoleFilter           string
	StatusFilter         string
	SortBy               string
	Page                 int
	PageSize             int
}

type petugasRepository struct {
	DB *gorm.DB
}

func NewPetugasRepository(db *gorm.DB) *petugasRepository {
	return &petugasRepository{DB: db}
}

func (r *petugasRepository) GetAll(params ParamsGetAllPetugas) ([]model.Petugas, pagination.Metadata, error) {
	var petugas []model.Petugas
	var totalRecords int64

	db := r.DB.Model(&model.Petugas{})

	if params.NameOrUsernameFilter != "" {
		searchQuery := "%" + params.NameOrUsernameFilter + "%"
		db = db.Where("nama_petugas ILIKE ? OR username_petugas ILIKE ?", searchQuery, searchQuery)
	}
	if params.RoleFilter != "" {
		db = db.Where("role = ?", params.RoleFilter)
	}
	if params.StatusFilter != "" {
		db = db.Where("status = ?", params.StatusFilter)
	}

	if err := db.Count(&totalRecords).Error; err != nil {
		return nil, pagination.Metadata{}, err
	}

	metadata := pagination.CalculateMetadata(int(totalRecords), params.Page, params.PageSize)

	sortWhiteList := map[string]string{
		"nama_asc":  "nama_petugas ASC",
		"nama_desc": "nama_petugas DESC",
	}

	orderByClause := "created_at DESC"
	if sort, ok := sortWhiteList[params.SortBy]; ok {
		orderByClause = sort
	}
	db = db.Order(orderByClause)

	db = db.Limit(metadata.PageSize).Offset((metadata.CurrentPage - 1) * metadata.PageSize)

	if err := db.Find(&petugas).Error; err != nil {
		return nil, pagination.Metadata{}, err
	}
	return petugas, metadata, nil
}

func (r *petugasRepository) GetById(id int) (model.Petugas, error) {
	var petugas model.Petugas
	result := r.DB.First(&petugas, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Petugas{}, ErrNotFound
		}
		return model.Petugas{}, result.Error
	}
	return petugas, nil
}

func (r *petugasRepository) GetByUsername(username string) (model.Petugas, error) {
	var petugas model.Petugas
	result := r.DB.Where("username_petugas = ?", username).First(&petugas)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Petugas{}, ErrNotFound
		}
		return model.Petugas{}, result.Error
	}
	return petugas, nil
}

func (r *petugasRepository) Create(petugas model.Petugas) (model.Petugas, error) {
	result := r.DB.Create(&petugas)
	return petugas, result.Error
}

func (r *petugasRepository) Update(id int, petugas model.Petugas) (model.Petugas, error) {
	petugas.ID = id

	result := r.DB.Model(&model.Petugas{}).Where("id_petugas = ?", id).Updates(&petugas)
	if result.Error != nil {
		return model.Petugas{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.Petugas{}, errors.New("data not found")
	}
	return r.GetById(id)
}

func (r *petugasRepository) Delete(id int) error {
	result := r.DB.Delete(&model.Petugas{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
