package db

type PaginationData struct {
	// PageNumber should start from 0
	PageNumber int32 `json:"page_number"`
	PageSize   int32 `json:"page_size"`
}

func (p *PaginationData) IsValid() bool {
	return p.PageNumber >= 0 && p.PageSize >= 1
}
