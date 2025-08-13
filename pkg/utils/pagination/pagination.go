package pagination

import "math"

// Metadata menampung informasi paginasi untuk respons JSON
type Metadata struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	TotalRecords int `json:"total_records"`
	TotalPages   int `json:"total_pages"`
}

// CalculateMetadata adalah fungsi helper untuk menghitung metadata
func CalculateMetadata(totalRecords, page, pageSize int) Metadata {
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
