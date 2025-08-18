package model

import (
	"database/sql"
	"time"
)

type Pemeriksaan struct {
	ID                 int            `json:"id,omitempty" gorm:"primaryKey;column:id_pemeriksaan"`
	AntrianID          int            `json:"antrian_id" gorm:"column:id_antrian;unique"`
	IcdID              sql.NullInt64  `json:"icd_id" gorm:"column:id_icd"`
	Nadi               sql.NullString `json:"nadi" gorm:"column:nadi"`
	TekananDarah       sql.NullString `json:"tekanan_darah" gorm:"column:tekanan_darah"`
	Suhu               sql.NullString `json:"suhu" gorm:"column:suhu"`
	BeratBadan         sql.NullString `json:"berat_badan" gorm:"column:berat_badan"`
	KeadaanUmum        sql.NullString `json:"keadaan_umum" gorm:"column:keadaan_umum"`
	Keluhan            sql.NullString `json:"keluhan" gorm:"column:keluhan"`
	RiwayatPenyakit    sql.NullString `json:"riwayat_penyakit" gorm:"column:riwayat_penyakit"`
	Keterangan         sql.NullString `json:"keterangan" gorm:"column:keterangan"`
	Tindakan           sql.NullString `json:"tindakan" gorm:"column:tindakan"`
	TanggalPemeriksaan string         `json:"tanggal_pemeriksaan" gorm:"column:tanggal_pemeriksaan"`
	CreatedAt          time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt          time.Time      `json:"updated_at" gorm:"column:updated_at"`

	Antrian Antrian `json:"antrian" gorm:"foreignKey:AntrianID"`
	Icd     Icd     `json:"icd" gorm:"foreignKey:IcdID"`
}

func (Pemeriksaan) TableName() string {
	return "pemeriksaan"
}

type CreatePemeriksaanRequest struct {
	AntrianID          int    `json:"antrian_id" binding:"required,gt=0"`
	IcdID              *int   `json:"icd_id,omitempty"`
	Nadi               string `json:"nadi,omitempty"`
	TekananDarah       string `json:"tekanan_darah,omitempty"`
	Suhu               string `json:"suhu,omitempty"`
	BeratBadan         string `json:"berat_badan,omitempty"`
	KeadaanUmum        string `json:"keadaan_umum,omitempty"`
	Keluhan            string `json:"keluhan,omitempty"`
	RiwayatPenyakit    string `json:"riwayat_penyakit,omitempty"`
	Keterangan         string `json:"keterangan,omitempty"`
	Tindakan           string `json:"tindakan,omitempty"`
	TanggalPemeriksaan string `json:"tanggal_pemeriksaan" binding:"required,datetime=2006-01-02"`
}

type UpdatePemeriksaanRequest struct {
	IcdID              *int   `json:"icd_id,omitempty"`
	Nadi               string `json:"nadi,omitempty"`
	TekananDarah       string `json:"tekanan_darah,omitempty"`
	Suhu               string `json:"suhu,omitempty"`
	BeratBadan         string `json:"berat_badan,omitempty"`
	KeadaanUmum        string `json:"keadaan_umum,omitempty"`
	Keluhan            string `json:"keluhan,omitempty"`
	RiwayatPenyakit    string `json:"riwayat_penyakit,omitempty"`
	Keterangan         string `json:"keterangan,omitempty"`
	Tindakan           string `json:"tindakan,omitempty"`
	TanggalPemeriksaan string `json:"tanggal_pemeriksaan" binding:"required,datetime=2006-01-02"`
}

type PemeriksaanDetail struct {
	ID                 int            `json:"id"`
	TanggalPemeriksaan string         `json:"tanggal_pemeriksaan"`
	Nadi               sql.NullString `json:"nadi"`
	TekananDarah       sql.NullString `json:"tekanan_darah"`
	Suhu               sql.NullString `json:"suhu"`
	BeratBadan         sql.NullString `json:"berat_badan"`
	Keluhan            sql.NullString `json:"keluhan"`
	Tindakan           sql.NullString `json:"tindakan"`
	Pasien             struct {
		ID           int    `json:"id"`
		Nama         string `json:"nama"`
		NoRekamMedis string `json:"no_rekam_medis"`
	} `json:"pasien"`
	Dokter struct {
		ID   int    `json:"id"`
		Nama string `json:"nama"`
	} `json:"dokter"`
	Poli struct {
		ID   int    `json:"id"`
		Nama string `json:"nama"`
	} `json:"poli"`
	Diagnosis struct {
		ID       sql.NullInt64  `json:"id"`
		Kode     sql.NullString `json:"kode"`
		Penyakit sql.NullString `json:"penyakit"`
	} `json:"diagnosis"`
}
