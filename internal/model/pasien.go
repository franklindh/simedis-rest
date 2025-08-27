package model

import (
	"database/sql"
	"time"
)

type Pasien struct {
	ID                        int            `json:"id,omitempty" gorm:"primaryKey;column:id_pasien"`
	NIK                       string         `json:"nik" gorm:"column:nik;unique"`
	NoRekamMedis              sql.NullString `json:"no_rekam_medis" gorm:"column:no_rekam_medis;unique"`
	NoKartuJaminan            sql.NullString `json:"no_kartu_jaminan" gorm:"column:no_kartu_jaminan"`
	UsernamePasien            string         `json:"username_pasien" gorm:"column:username_pasien;unique"`
	NoTeleponPasien           sql.NullString `json:"no_telepon_pasien" gorm:"column:no_telepon_pasien"`
	NamaPasien                string         `json:"nama_pasien" gorm:"column:nama_pasien"`
	AlamatPasien              string         `json:"alamat_pasien" gorm:"column:alamat_pasien"`
	TempatLahirPasien         string         `json:"tempat_lahir_pasien" gorm:"column:tempat_lahir_pasien"`
	TanggalLahirPasien        time.Time      `json:"tanggal_lahir_pasien" gorm:"column:tanggal_lahir_pasien"`
	JKPasien                  string         `json:"jk_pasien" gorm:"column:jk_pasien"`
	StatusPernikahan          string         `json:"status_pernikahan" gorm:"column:status_pernikahan"`
	NamaKeluargaTerdekat      sql.NullString `json:"nama_keluarga_terdekat" gorm:"column:nama_keluarga_terdekat"`
	NoTeleponKeluargaTerdekat sql.NullString `json:"no_telepon_keluarga_terdekat" gorm:"column:no_telepon_keluarga_terdekat"`
	Password                  string         `json:"-" gorm:"column:password"`
	CreatedAt                 time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt                 time.Time      `json:"updated_at" gorm:"column:updated_at"`
}

func (Pasien) TableName() string {
	return "pasien"
}

type CreatePasienRequest struct {
	NIK                       string `json:"nik" binding:"required,numeric,len=16"`
	NoKartuJaminan            string `json:"no_kartu_jaminan,omitempty"`
	UsernamePasien            string `json:"username_pasien,omitempty" binding:"sanitize"`
	Password                  string `json:"password,omitempty" binding:"omitempty,min=8"`
	NoTeleponPasien           string `json:"no_telepon_pasien,omitempty"`
	NamaPasien                string `json:"nama_pasien" binding:"required,sanitize"`
	AlamatPasien              string `json:"alamat_pasien" binding:"required,sanitize"`
	TempatLahirPasien         string `json:"tempat_lahir_pasien" binding:"required,sanitize"`
	TanggalLahirPasien        string `json:"tanggal_lahir_pasien" binding:"required,datetime=2006-01-02"`
	JKPasien                  string `json:"jk_pasien" binding:"required,oneof=L P"`
	StatusPernikahan          string `json:"status_pernikahan" binding:"required,oneof='Belum Menikah' Menikah 'Cerai Hidup' 'Cerai Mati'"`
	NamaKeluargaTerdekat      string `json:"nama_keluarga_terdekat,omitempty" binding:"sanitize"`
	NoTeleponKeluargaTerdekat string `json:"no_telepon_keluarga_terdekat,omitempty"`
}

type UpdatePasienRequest struct {
	NIK                       string `json:"nik" binding:"required,numeric,len=16"`
	NoKartuJaminan            string `json:"no_kartu_jaminan,omitempty"`
	UsernamePasien            string `json:"username_pasien,omitempty" binding:"sanitize"`
	NoTeleponPasien           string `json:"no_telepon_pasien,omitempty"`
	NamaPasien                string `json:"nama_pasien" binding:"required,sanitize"`
	AlamatPasien              string `json:"alamat_pasien" binding:"required,sanitize"`
	TempatLahirPasien         string `json:"tempat_lahir_pasien" binding:"required,sanitize"`
	TanggalLahirPasien        string `json:"tanggal_lahir_pasien" binding:"required,datetime=2006-01-02"`
	JKPasien                  string `json:"jk_pasien" binding:"required,oneof=L P"`
	StatusPernikahan          string `json:"status_pernikahan" binding:"required,oneof='Belum Menikah' Menikah 'Cerai Hidup' 'Cerai Mati'"`
	NamaKeluargaTerdekat      string `json:"nama_keluarga_terdekat,omitempty" binding:"sanitize"`
	NoTeleponKeluargaTerdekat string `json:"no_telepon_keluarga_terdekat,omitempty"`
}

func (req *CreatePasienRequest) ToModel(username, hashedPassword, noRekamMedis string) Pasien {
	parsedDate, _ := time.Parse("2006-01-02", req.TanggalLahirPasien)
	return Pasien{
		NIK:                       req.NIK,
		NoKartuJaminan:            sql.NullString{String: req.NoKartuJaminan, Valid: req.NoKartuJaminan != ""},
		UsernamePasien:            username,
		Password:                  hashedPassword,
		NoTeleponPasien:           sql.NullString{String: req.NoTeleponPasien, Valid: req.NoTeleponPasien != ""},
		NamaPasien:                req.NamaPasien,
		AlamatPasien:              req.AlamatPasien,
		TempatLahirPasien:         req.TempatLahirPasien,
		TanggalLahirPasien:        parsedDate,
		JKPasien:                  req.JKPasien,
		StatusPernikahan:          req.StatusPernikahan,
		NamaKeluargaTerdekat:      sql.NullString{String: req.NamaKeluargaTerdekat, Valid: req.NamaKeluargaTerdekat != ""},
		NoTeleponKeluargaTerdekat: sql.NullString{String: req.NoTeleponKeluargaTerdekat, Valid: req.NoTeleponKeluargaTerdekat != ""},
		NoRekamMedis:              sql.NullString{String: noRekamMedis, Valid: true},
	}
}

func (req *UpdatePasienRequest) ToModel() Pasien {
	parsedDate, _ := time.Parse("2006-01-02", req.TanggalLahirPasien)
	return Pasien{
		NIK:                       req.NIK,
		NoKartuJaminan:            sql.NullString{String: req.NoKartuJaminan, Valid: req.NoKartuJaminan != ""},
		UsernamePasien:            req.UsernamePasien,
		NoTeleponPasien:           sql.NullString{String: req.NoTeleponPasien, Valid: req.NoTeleponPasien != ""},
		NamaPasien:                req.NamaPasien,
		AlamatPasien:              req.AlamatPasien,
		TempatLahirPasien:         req.TempatLahirPasien,
		TanggalLahirPasien:        parsedDate,
		JKPasien:                  req.JKPasien,
		StatusPernikahan:          req.StatusPernikahan,
		NamaKeluargaTerdekat:      sql.NullString{String: req.NamaKeluargaTerdekat, Valid: req.NamaKeluargaTerdekat != ""},
		NoTeleponKeluargaTerdekat: sql.NullString{String: req.NoTeleponKeluargaTerdekat, Valid: req.NoTeleponKeluargaTerdekat != ""},
	}
}

type PasienResponse struct {
	ID                        int       `json:"id"`
	NIK                       string    `json:"nik"`
	NoRekamMedis              string    `json:"no_rekam_medis,omitempty"`
	NoKartuJaminan            string    `json:"no_kartu_jaminan,omitempty"`
	UsernamePasien            string    `json:"username_pasien"`
	NoTeleponPasien           string    `json:"no_telepon_pasien,omitempty"`
	NamaPasien                string    `json:"nama_pasien"`
	AlamatPasien              string    `json:"alamat_pasien"`
	TempatLahirPasien         string    `json:"tempat_lahir_pasien"`
	TanggalLahirPasien        time.Time `json:"tanggal_lahir_pasien"`
	JKPasien                  string    `json:"jk_pasien"`
	StatusPernikahan          string    `json:"status_pernikahan"`
	NamaKeluargaTerdekat      string    `json:"nama_keluarga_terdekat,omitempty"`
	NoTeleponKeluargaTerdekat string    `json:"no_telepon_keluarga_terdekat,omitempty"`
	CreatedAt                 time.Time `json:"created_at"`
}

func ToPasienResponse(p Pasien) PasienResponse {
	return PasienResponse{
		ID:                        p.ID,
		NIK:                       p.NIK,
		NoRekamMedis:              p.NoRekamMedis.String,
		NoKartuJaminan:            p.NoKartuJaminan.String,
		UsernamePasien:            p.UsernamePasien,
		NoTeleponPasien:           p.NoTeleponPasien.String,
		NamaPasien:                p.NamaPasien,
		AlamatPasien:              p.AlamatPasien,
		TempatLahirPasien:         p.TempatLahirPasien,
		TanggalLahirPasien:        p.TanggalLahirPasien,
		JKPasien:                  p.JKPasien,
		StatusPernikahan:          p.StatusPernikahan,
		NamaKeluargaTerdekat:      p.NamaKeluargaTerdekat.String,
		NoTeleponKeluargaTerdekat: p.NoTeleponKeluargaTerdekat.String,
		CreatedAt:                 p.CreatedAt,
	}
}

func ToPasienResponseList(pasiens []Pasien) []PasienResponse {
	var responses []PasienResponse
	for _, p := range pasiens {
		responses = append(responses, ToPasienResponse(p))
	}
	return responses
}
