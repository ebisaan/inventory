package domain

const (
	DefaultPageSize = 50
	MaxPageSize     = 100
)

type Filter struct {
	Page     int
	PageSize int
}

func ProcessFilter(filter Filter) Filter {
	if filter.PageSize == 0 {
		filter.PageSize = DefaultPageSize
	}

	if filter.PageSize > MaxPageSize {
		filter.PageSize = MaxPageSize
	}

	if filter.Page == 0 {
		filter.Page = 1
	}

	return filter
}

func (f Filter) Limit() int {
	return f.PageSize
}

func (f Filter) Offset() int64 {
	return int64(f.PageSize) * int64(f.Page-1)
}
