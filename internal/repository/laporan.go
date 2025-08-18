package repository

import (
	"github.com/franklindh/simedis-api/internal/model"
	"gorm.io/gorm"
)

type LaporanRepository struct {
	DB *gorm.DB
}

func NewLaporanRepository(db *gorm.DB) *LaporanRepository {
	return &LaporanRepository{DB: db}
}

func (r *LaporanRepository) GetLaporanKunjunganPerPoli(startDate, endDate string) ([]model.LaporanKunjunganPoli, error) {
	var results []model.LaporanKunjunganPoli

	err := r.DB.Table("antrian").
		Select("poli.nama_poli, count(antrian.id_antrian) as jumlah_kunjungan").
		Joins("join jadwal on antrian.id_jadwal = jadwal.id_jadwal").
		Joins("join poli on jadwal.id_poli = poli.id_poli").
		Where("jadwal.tanggal_praktik BETWEEN ? AND ?", startDate, endDate).
		Group("poli.nama_poli").
		Scan(&results).Error

	return results, err
}

func (r *LaporanRepository) GetLaporanPenyakitTeratas(startDate, endDate string, limit int) ([]model.LaporanPenyakitTeratas, error) {
	var results []model.LaporanPenyakitTeratas

	err := r.DB.Table("pemeriksaan").
		Select("icd.kode_icd, icd.nama_penyakit, count(pemeriksaan.id_pemeriksaan) as jumlah_kasus").
		Joins("join icd on pemeriksaan.id_icd = icd.id_icd").
		Where("pemeriksaan.tanggal_pemeriksaan BETWEEN ? AND ?", startDate, endDate).
		Where("pemeriksaan.id_icd IS NOT NULL").
		Group("icd.kode_icd, icd.nama_penyakit").
		Order("jumlah_kasus DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}
