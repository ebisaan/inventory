package domain

type Filter struct {
	Page     int
	PageSize int
}

func (f Filter) Limit() int {
	return f.PageSize
}

func (f Filter) Offset() int64 {
	return int64(f.PageSize) * int64(f.Page-1)
}
