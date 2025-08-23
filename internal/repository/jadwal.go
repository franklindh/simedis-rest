package repository

import (
	"errors"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/pkg/utils/pagination"
	"gorm.io/gorm"
)

type ParamsGetAllJadwal struct {
	PetugasIDFilter int
	PoliIDFilter    int
	StartDateFilter string
	EndDateFilter   string
	SortBy          string
	Page            int
	PageSize        int
}

type JadwalRepository struct {
	DB *gorm.DB
}

func NewJadwalRepository(db *gorm.DB) *JadwalRepository {
	return &JadwalRepository{DB: db}
}

func (r *JadwalRepository) Create(jadwal model.Jadwal) (model.Jadwal, error) {
	result := r.DB.Create(&jadwal)
	return jadwal, result.Error
}

func (r *JadwalRepository) GetAll(params ParamsGetAllJadwal) ([]model.Jadwal, pagination.Metadata, error) {
	var jadwal []model.Jadwal
	var totalRecords int64

	db := r.DB.Model(&model.Jadwal{}).Preload("Petugas").Preload("Poli")

	if params.PoliIDFilter > 0 {
		db = db.Where("id_poli = ?", params.PoliIDFilter)
	}
	if params.PetugasIDFilter > 0 {
		db = db.Where("id_petugas = ?", params.PetugasIDFilter)
	}
	if params.StartDateFilter != "" {
		db = db.Where("tanggal_praktik >= ?", params.StartDateFilter)
	}
	if params.EndDateFilter != "" {
		db = db.Where("tanggal_praktik <= ?", params.EndDateFilter)
	}

	if err := db.Count(&totalRecords).Error; err != nil {
		return nil, pagination.Metadata{}, err
	}

	metadata := pagination.CalculateMetadata(int(totalRecords), params.Page, params.PageSize)

	sortWhiteList := map[string]string{
		"tanggal_asc":  "tanggal_praktik ASC",
		"tanggal_desc": "tanggal_praktik DESC",
	}
	orderByClause := "tanggal_praktik DESC"
	if sort, ok := sortWhiteList[params.SortBy]; ok {
		orderByClause = sort
	}
	db = db.Order(orderByClause)

	db = db.Limit(metadata.PageSize).Offset((metadata.CurrentPage - 1) * metadata.PageSize)

	if err := db.Find(&jadwal).Error; err != nil {
		return nil, pagination.Metadata{}, err
	}
	return jadwal, metadata, nil
}

func (r *JadwalRepository) GetById(id int) (model.Jadwal, error) {
	var jadwal model.Jadwal
	result := r.DB.Preload("Petugas").Preload("Poli").First(&jadwal, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Jadwal{}, ErrNotFound
		}
		return model.Jadwal{}, result.Error
	}
	return jadwal, nil
}

func (r *JadwalRepository) Update(id int, jadwal model.Jadwal) (model.Jadwal, error) {
	jadwal.ID = id
	result := r.DB.Model(&model.Jadwal{}).Where("id_jadwal = ?", id).Updates(&jadwal)
	if result.Error != nil {
		return model.Jadwal{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.Jadwal{}, ErrNotFound
	}
	return r.GetById(id)
}

func (r *JadwalRepository) Delete(id int) error {
	result := r.DB.Delete(&model.Jadwal{}, id)
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}
