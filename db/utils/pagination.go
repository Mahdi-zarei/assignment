package db

// DefaultPagination can be used to retrieve the first 30 records
var DefaultPagination = PaginationData{
	PageNumber: 0,
	PageSize:   30,
}

type PaginationData struct {
	// PageNumber should start from 0
	PageNumber int32 `json:"page_number"`
	PageSize   int32 `json:"page_size"`
}

func (p *PaginationData) IsValid() bool {
	return p.PageNumber >= 0 && p.PageSize >= 1
}

func (p *PaginationData) GetLimitAndOffset() (limit int32, offset int32) {
	limit = p.PageSize
	offset = p.PageNumber * p.PageSize
	return
}
