package repository

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/franklindh/simedis-api/internal/model"
)

type PetugasRepository struct {
	DB *sql.DB
}

type ParamsGetAllParams struct {
	NameFilter string
	SortBy     string
	Page       int
	PageSize   int
}

type Metadata struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	TotalRecords int `json:"total_records"`
	TotalPages   int `json:"total_pages"`
}

func NewPetugasRepository(db *sql.DB) *PetugasRepository {
	return &PetugasRepository{DB: db}
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}
	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		TotalRecords: totalRecords,
		TotalPages:   int(math.Ceil(float64(totalRecords) / float64(pageSize))),
	}
}

// func (r *PetugasRepository) GetAll() ([]model.Petugas, error) {
// 	query := `SELECT id_petugas, id_poli, username_petugas, nama_petugas, status, role, created_at, updated_at FROM petugas WHERE deleted_at IS NULL`

// 	rows, err := r.DB.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var petugas []model.Petugas
// 	for rows.Next() {
// 		var p model.Petugas
// 		if err := rows.Scan(&p.ID, &p.PoliID, &p.Username, &p.Name, &p.Status, &p.Role, &p.CreatedAt, &p.UpdatedAt); err != nil {
// 			return nil, err
// 		}
// 		petugas = append(petugas, p)
// 	}

//		return petugas, nil
//	}
func (r *PetugasRepository) GetAll(params ParamsGetAllParams) ([]model.Petugas, Metadata, error) {
	query := ` FROM petugas WHERE deleted_at IS NULL`
	dataQuery := `SELECT id_petugas, id_poli, username_petugas, nama_petugas, status, role, created_at, updated_at` + query
	countQuery := `SELECT COUNT(*)` + query

	whereClause := ""
	args := []any{}
	if params.NameFilter != "" {
		whereClause = " AND (nama_petugas ILIKE $1 OR username_petugas ILIKE $1)"
		args = append(args, "%"+params.NameFilter+"%")
	}

	var totalRecords int
	err := r.DB.QueryRow(countQuery+whereClause, args...).Scan(&totalRecords)
	if err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, params.Page, params.PageSize)

	sortWhitelist := map[string]string{
		"name_asc":  "nama_petugas ASC",
		"name_desc": "nama_petugas DESC",
	}
	orderByClause := " ORDER BY created_at DESC"
	if sort, ok := sortWhitelist[params.SortBy]; ok {
		orderByClause = " ORDER BY " + sort
	}

	paginationClause := fmt.Sprintf(" LIMIT %d OFFSET %d", metadata.PageSize, (metadata.CurrentPage-1)*metadata.PageSize)

	finalQuery := dataQuery + whereClause + orderByClause + paginationClause

	rows, err := r.DB.Query(finalQuery, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	var allPetugas []model.Petugas
	for rows.Next() {
		var p model.Petugas
		if err := rows.Scan(&p.ID, &p.PoliID, &p.Username, &p.Name, &p.Status, &p.Role, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, Metadata{}, err
		}
		allPetugas = append(allPetugas, p)
	}

	return allPetugas, metadata, nil
}

func (r *PetugasRepository) GetByID(id int) (model.Petugas, error) {
	query := `SELECT id_petugas, id_poli, username_petugas, nama_petugas, status, role, created_at, updated_at FROM petugas WHERE id_petugas = $1 AND deleted_at IS NULL`

	var p model.Petugas
	err := r.DB.QueryRow(query, id).Scan(&p.ID, &p.PoliID, &p.Username, &p.Name, &p.Status, &p.Role, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Petugas{}, ErrNotFound
		}
		return model.Petugas{}, err
	}

	return p, nil
}

func (r *PetugasRepository) Create(petugas model.Petugas) (model.Petugas, error) {
	query := `
		INSERT INTO petugas (id_poli, username_petugas, nama_petugas, status, role, password)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id_petugas, created_at, updated_at`

	err := r.DB.QueryRow(query, petugas.PoliID, petugas.Username, petugas.Name, petugas.Status, petugas.Role, petugas.Password).Scan(&petugas.ID, &petugas.CreatedAt, &petugas.UpdatedAt)

	if err != nil {
		return model.Petugas{}, err
	}

	return petugas, nil
}

func (r *PetugasRepository) Update(id int, petugas model.Petugas) (model.Petugas, error) {
	query := `UPDATE petugas SET id_poli = $1, username_petugas = $2, nama_petugas = $3, status = $4, role = $5, updated_at = NOW() WHERE id_petugas = $6 AND deleted_at IS NULL RETURNING id_petugas, id_poli, username_petugas, nama_petugas, status, role, created_at, updated_at`

	var updatedPetugas model.Petugas
	err := r.DB.QueryRow(query, petugas.PoliID, petugas.Username, petugas.Name, petugas.Status, petugas.Role, id).Scan(&updatedPetugas.ID, &updatedPetugas.PoliID, &updatedPetugas.Username, &updatedPetugas.Name, &updatedPetugas.Status, &updatedPetugas.Role, &updatedPetugas.CreatedAt, &updatedPetugas.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.Petugas{}, ErrNotFound
		}
		return model.Petugas{}, err
	}

	return updatedPetugas, nil
}

func (r *PetugasRepository) Delete(id int) error {
	query := `UPDATE petugas SET deleted_at = NOW() WHERE id_petugas = $1 AND deleted_at IS NULL`

	res, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
