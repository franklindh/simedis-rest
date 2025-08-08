package repository

import (
	"database/sql"
	"errors"

	"github.com/franklindh/simedis-api/internal/model"
)

var ErrNotFound = errors.New("data not found")

type PoliRepository struct {
	DB *sql.DB
}

func NewPoliRepository(db *sql.DB) *PoliRepository {
	return &PoliRepository{DB: db}
}

func (r *PoliRepository) GetAll() ([]model.Poli, error) {
	query := `SELECT id_poli, nama_poli, status_poli, created_at, updated_at FROM poli WHERE deleted_at IS NULL`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var poli []model.Poli
	for rows.Next() {
		var p model.Poli
		if err := rows.Scan(&p.ID, &p.Name, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		poli = append(poli, p)
	}
	return poli, nil
}

func (r *PoliRepository) GetByID(id int) (model.Poli, error) {
	query := `SELECT id_poli, nama_poli, status_poli, created_at, updated_at FROM poli WHERE id_poli = $1 AND deleted_at IS NULL`

	var p model.Poli
	err := r.DB.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Poli{}, ErrNotFound
		}
		return model.Poli{}, err
	}
	return p, nil
}

func (r *PoliRepository) Create(poli model.Poli) (model.Poli, error) {
	query := `INSERT INTO poli (nama_poli, status_poli) VALUES ($1, $2) RETURNING id_poli, created_at, updated_at`

	err := r.DB.QueryRow(query, poli.Name, poli.Status).Scan(&poli.ID, &poli.CreatedAt, &poli.UpdatedAt)
	if err != nil {
		return model.Poli{}, err
	}

	return poli, nil
}

// func (r *PoliRepository) Update(id int, updates map[string]any) (model.Poli, error) {
// 	jsonToDBMap := map[string]string{
// 		"nama":   "nama_poli",
// 		"status": "status_poli",
// 	}

// 	setClauses := make([]string, 0)
// 	args := make([]any, 0)
// 	argId := 1

// 	for key, value := range updates {
// 		dbColumn, ok := jsonToDBMap[key]
// 		if !ok {
// 			continue
// 		}
// 		setClauses = append(setClauses, fmt.Sprintf("%s =  $%d", dbColumn, argId))
// 		args = append(args, value)
// 		argId++
// 	}

// 	if len(setClauses) == 0 {
// 		return model.Poli{}, fmt.Errorf("no fields to update")
// 	}

// 	args = append(args, id)

// 	query := fmt.Sprintf("UPDATE poli SET %s WHERE id_poli = $%d RETURNING id_poli, nama_poli, status_poli, created_at, updated_at", strings.Join(setClauses, ", "), argId)

// 	var updatedPoli model.Poli

// 	err := r.DB.QueryRow(query, args...).Scan(&updatedPoli.ID, &updatedPoli.Name, &updatedPoli.Status, &updatedPoli.CreatedAt, &updatedPoli.UpdatedAt)

// 	if err != nil {
// 		return model.Poli{}, err
// 	}

// 	return updatedPoli, nil
// }

func (r *PoliRepository) Update(id int, poli model.Poli) (model.Poli, error) {
	query := `UPDATE poli SET nama_poli = $1, status_poli = $2, updated_at = NOW() WHERE id_poli = $3 RETURNING id_poli, nama_poli, status_poli, created_at, updated_at`

	var updatePoli model.Poli

	err := r.DB.QueryRow(query, poli.Name, poli.Status, id).Scan(&updatePoli.ID, &updatePoli.Name, &updatePoli.Status, &updatePoli.CreatedAt, &updatePoli.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.Poli{}, ErrNotFound
		}
		return model.Poli{}, err
	}

	return updatePoli, nil
}

func (r *PoliRepository) Delete(id int) error {
	query := `UPDATE poli SET deleted_at = NOW() WHERE id_poli = $1 AND deleted_at IS NULL`

	res, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffedcted, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffedcted == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *PoliRepository) FindByNameIncludingDeleted(name string) (model.Poli, error) {
	query := `SELECT id_poli, nama_poli, status_poli, created_at, updated_at, deleted_at FROM poli WHERE nama_poli = $1`

	var p model.Poli
	err := r.DB.QueryRow(query, name).Scan(&p.ID, &p.Name, &p.Status, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt)
	if err != nil {
		return model.Poli{}, err
	}
	return p, nil
}

func (r *PoliRepository) Restore(id int) (model.Poli, error) {
	query := `UPDATE poli SET deleted_at = NULL, updated_at = NOW() WHERE id_poli = $1 AND deleted_at IS NOT NULL RETURNING id_poli, nama_poli, status_poli, created_at, updated_at`

	var restoredPoli model.Poli
	err := r.DB.QueryRow(query, id).Scan(&restoredPoli.ID, &restoredPoli.Name, &restoredPoli.Status, &restoredPoli.CreatedAt, &restoredPoli.UpdatedAt)

	return restoredPoli, err
}
