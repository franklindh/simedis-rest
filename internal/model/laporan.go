package model

type LaporanKunjunganPoli struct {
	NamaPoli        string `json:"nama_poli"`
	JumlahKunjungan int    `json:"jumlah_kunjungan"`
}

type LaporanPenyakitTeratas struct {
	KodeIcd      string `json:"kode_icd"`
	NamaPenyakit string `json:"nama_penyakit"`
	JumlahKasus  int    `json:"jumlah_kasus"`
}
