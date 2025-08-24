package repository

import (
	"errors"
	"time"

	"github.com/franklindh/simedis-api/internal/model"
	"github.com/franklindh/simedis-api/pkg/utils/pagination"
	"gorm.io/gorm"
)

type ParamsGetAllAntrian struct {
	StatusFilter  string
	TanggalFilter string
	PoliIDFilter  int
	SortBy        string
	Page          int
	PageSize      int
}

type AntrianRepository struct {
	DB *gorm.DB
}

func NewAntrianRepository(db *gorm.DB) *AntrianRepository {
	return &AntrianRepository{DB: db}
}

func (r *AntrianRepository) Create(antrian model.Antrian) (model.Antrian, error) {
	result := r.DB.Create(&antrian)
	if result.Error != nil {
		return model.Antrian{}, result.Error
	}

	return r.GetByID(antrian.ID)
}

func (r *AntrianRepository) GetAll(params ParamsGetAllAntrian) ([]model.Antrian, pagination.Metadata, error) {
	var antrian []model.Antrian
	var totalRecords int64

	db := r.DB.Model(&model.Antrian{}).Preload("Pasien").Preload("Jadwal.Poli").Preload("Jadwal.Petugas")

	if params.StatusFilter != "" {
		db = db.Where("status = ?", params.StatusFilter)
	}
	if params.TanggalFilter != "" || params.PoliIDFilter > 0 {
		db = db.Joins("JOIN jadwal ON antrian.id_jadwal = jadwal.id_jadwal")
		if params.TanggalFilter != "" {
			db = db.Where("jadwal.tanggal_praktik = ?", params.TanggalFilter)
		}
		if params.PoliIDFilter > 0 {
			db = db.Where("jadwal.id_poli = ?", params.PoliIDFilter)
		}
	}

	if err := db.Count(&totalRecords).Error; err != nil {
		return nil, pagination.Metadata{}, err
	}

	metadata := pagination.CalculateMetadata(int(totalRecords), params.Page, params.PageSize)

	db = db.Order("created_at ASC")

	db = db.Limit(metadata.PageSize).Offset((metadata.CurrentPage - 1) * metadata.PageSize)

	if err := db.Find(&antrian).Error; err != nil {
		return nil, pagination.Metadata{}, err
	}

	return antrian, metadata, nil
}

func (r *AntrianRepository) GetByID(id int) (model.Antrian, error) {
	var antrian model.Antrian
	result := r.DB.Preload("Pasien").Preload("Jadwal.Poli").Preload("Jadwal.Poli").Preload("Jadwal.Petugas").First(&antrian, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return model.Antrian{}, ErrNotFound
		}
		return model.Antrian{}, result.Error
	}
	return antrian, nil
}

func (r *AntrianRepository) Update(id int, antrian model.Antrian) (model.Antrian, error) {
	result := r.DB.Model(&antrian).Where("id_antrian = ?", id).Updates(antrian)
	if result.Error != nil {
		return model.Antrian{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.Antrian{}, ErrNotFound
	}
	return r.GetByID(id)
}

func (r *AntrianRepository) Delete(id int) error {
	result := r.DB.Delete(&model.Antrian{}, id)
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return result.Error
}

func (r *AntrianRepository) CheckAntrian(pasienID, jadwalID int) (bool, error) {
	var count int64

	result := r.DB.Model(&model.Antrian{}).
		Where("id_pasien = ?", pasienID).
		Where("id_jadwal = ?", jadwalID).
		Where("status IN ?", []string{"Menunggu", "Menunggu Diagnosis"}).
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (r *AntrianRepository) CheckForOverlappingAntrian(pasienID int, tanggal time.Time, waktuMulai time.Time, waktuSelesai time.Time) (bool, error) {
	var count int64

	result := r.DB.Model(&model.Antrian{}).
		Joins("JOIN jadwal ON antrian.id_jadwal = jadwal.id_jadwal").
		Where("antrian.id_pasien = ?", pasienID).
		Where("jadwal.tanggal_praktik = ?", tanggal).
		Where("jadwal.waktu_selesai > ?", waktuMulai).
		Where("jadwal.waktu_mulai < ?", waktuSelesai).
		Where("antrian.status IN ?", []string{"Menunggu", "Menunggu Diagnosis"}).
		Count(&count)

	if result.Error != nil {
		return false, result.Error
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (r *AntrianRepository) CountTodayByJadwal(jadwalID int) (int64, error) {
	var count int64

	result := r.DB.Model(&model.Antrian{}).
		Where("id_jadwal = ?", jadwalID).
		Where("CAST(created_at AS DATE) = CURRENT_DATE").
		Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}
