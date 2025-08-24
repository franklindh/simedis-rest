package model

import (
	"database/sql"
	"time"
)

type Pemeriksaan struct {
	ID              int            `json:"id,omitempty" gorm:"primaryKey;column:id_pemeriksaan"`
	AntrianID       int            `json:"antrian_id" gorm:"column:id_antrian;unique"`
	IcdID           sql.NullInt64  `json:"icd_id" gorm:"column:id_icd"`
	Nadi            sql.NullString `json:"nadi" gorm:"column:nadi"`
	TekananDarah    sql.NullString `json:"tekanan_darah" gorm:"column:tekanan_darah"`
	Suhu            sql.NullString `json:"suhu" gorm:"column:suhu"`
	BeratBadan      sql.NullString `json:"berat_badan" gorm:"column:berat_badan"`
	KeadaanUmum     sql.NullString `json:"keadaan_umum" gorm:"column:keadaan_umum"`
	Keluhan         sql.NullString `json:"keluhan" gorm:"column:keluhan"`
	RiwayatPenyakit sql.NullString `json:"riwayat_penyakit" gorm:"column:riwayat_penyakit"`
	Keterangan      sql.NullString `json:"keterangan" gorm:"column:keterangan"`
	Tindakan        sql.NullString `json:"tindakan" gorm:"column:tindakan"`

	TanggalPemeriksaan time.Time `json:"tanggal_pemeriksaan" gorm:"column:tanggal_pemeriksaan"`
	CreatedAt          time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"column:updated_at"`

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

func (req *CreatePemeriksaanRequest) ToModel() Pemeriksaan {

	parsedDate, _ := time.Parse("2006-01-02", req.TanggalPemeriksaan)

	pemeriksaan := Pemeriksaan{
		AntrianID:          req.AntrianID,
		Nadi:               sql.NullString{String: req.Nadi, Valid: req.Nadi != ""},
		TekananDarah:       sql.NullString{String: req.TekananDarah, Valid: req.TekananDarah != ""},
		Suhu:               sql.NullString{String: req.Suhu, Valid: req.Suhu != ""},
		BeratBadan:         sql.NullString{String: req.BeratBadan, Valid: req.BeratBadan != ""},
		KeadaanUmum:        sql.NullString{String: req.KeadaanUmum, Valid: req.KeadaanUmum != ""},
		Keluhan:            sql.NullString{String: req.Keluhan, Valid: req.Keluhan != ""},
		RiwayatPenyakit:    sql.NullString{String: req.RiwayatPenyakit, Valid: req.RiwayatPenyakit != ""},
		Keterangan:         sql.NullString{String: req.Keterangan, Valid: req.Keterangan != ""},
		Tindakan:           sql.NullString{String: req.Tindakan, Valid: req.Tindakan != ""},
		TanggalPemeriksaan: parsedDate,
	}
	if req.IcdID != nil {
		pemeriksaan.IcdID = sql.NullInt64{Int64: int64(*req.IcdID), Valid: true}
	}
	return pemeriksaan
}

func (req *UpdatePemeriksaanRequest) ToModel() Pemeriksaan {

	parsedDate, _ := time.Parse("2006-01-02", req.TanggalPemeriksaan)

	pemeriksaan := Pemeriksaan{
		Nadi:               sql.NullString{String: req.Nadi, Valid: req.Nadi != ""},
		TekananDarah:       sql.NullString{String: req.TekananDarah, Valid: req.TekananDarah != ""},
		Suhu:               sql.NullString{String: req.Suhu, Valid: req.Suhu != ""},
		BeratBadan:         sql.NullString{String: req.BeratBadan, Valid: req.BeratBadan != ""},
		KeadaanUmum:        sql.NullString{String: req.KeadaanUmum, Valid: req.KeadaanUmum != ""},
		Keluhan:            sql.NullString{String: req.Keluhan, Valid: req.Keluhan != ""},
		RiwayatPenyakit:    sql.NullString{String: req.RiwayatPenyakit, Valid: req.RiwayatPenyakit != ""},
		Keterangan:         sql.NullString{String: req.Keterangan, Valid: req.Keterangan != ""},
		Tindakan:           sql.NullString{String: req.Tindakan, Valid: req.Tindakan != ""},
		TanggalPemeriksaan: parsedDate,
	}
	if req.IcdID != nil {
		pemeriksaan.IcdID = sql.NullInt64{Int64: int64(*req.IcdID), Valid: true}
	}
	return pemeriksaan
}

type PemeriksaanResponse struct {
	ID                 int           `json:"id"`
	TanggalPemeriksaan time.Time     `json:"tanggal_pemeriksaan"`
	Nadi               string        `json:"nadi,omitempty"`
	TekananDarah       string        `json:"tekanan_darah,omitempty"`
	Suhu               string        `json:"suhu,omitempty"`
	BeratBadan         string        `json:"berat_badan,omitempty"`
	Keluhan            string        `json:"keluhan,omitempty"`
	Tindakan           string        `json:"tindakan,omitempty"`
	Pasien             PasienInfo    `json:"pasien"`
	Dokter             PetugasInfo   `json:"dokter"`
	Poli               PoliInfo      `json:"poli"`
	Diagnosis          DiagnosisInfo `json:"diagnosis"`
}

type PasienInfo struct {
	ID           int    `json:"id"`
	Nama         string `json:"nama"`
	NoRekamMedis string `json:"no_rekam_medis"`
}

type PetugasInfo struct {
	ID   int    `json:"id"`
	Nama string `json:"nama"`
}

type PoliInfo struct {
	ID   int    `json:"id"`
	Nama string `json:"nama"`
}

type DiagnosisInfo struct {
	ID       *int64  `json:"id,omitempty"`
	Kode     *string `json:"kode,omitempty"`
	Penyakit *string `json:"penyakit,omitempty"`
}

func ToPemeriksaanResponse(p Pemeriksaan) PemeriksaanResponse {
	resp := PemeriksaanResponse{
		ID:                 p.ID,
		TanggalPemeriksaan: p.TanggalPemeriksaan,
		Nadi:               p.Nadi.String,
		TekananDarah:       p.TekananDarah.String,
		Suhu:               p.Suhu.String,
		BeratBadan:         p.BeratBadan.String,
		Keluhan:            p.Keluhan.String,
		Tindakan:           p.Tindakan.String,
		Pasien: PasienInfo{
			ID:           p.Antrian.Pasien.ID,
			Nama:         p.Antrian.Pasien.NamaPasien,
			NoRekamMedis: p.Antrian.Pasien.NoRekamMedis.String,
		},
		Dokter: PetugasInfo{
			ID:   p.Antrian.Jadwal.Petugas.ID,
			Nama: p.Antrian.Jadwal.Petugas.Nama,
		},
		Poli: PoliInfo{
			ID:   p.Antrian.Jadwal.Poli.ID,
			Nama: p.Antrian.Jadwal.Poli.Nama,
		},
	}

	if p.IcdID.Valid {

		id := p.IcdID.Int64
		kode := p.Icd.KodeIcd
		penyakit := p.Icd.NamaPenyakit

		resp.Diagnosis.ID = &id
		resp.Diagnosis.Kode = &kode
		resp.Diagnosis.Penyakit = &penyakit
	}
	return resp
}

func ToPemeriksaanResponseList(pemeriksaanList []Pemeriksaan) []PemeriksaanResponse {
	var responses []PemeriksaanResponse
	for _, p := range pemeriksaanList {
		responses = append(responses, ToPemeriksaanResponse(p))
	}
	return responses
}
