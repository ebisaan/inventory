package domain

import (
	"math"
)

type Metadata struct {
	CurrentPage  int
	FirstPage    int
	LastPage     int
	PageSize     int
	TotalRecords int64
}

func MakeMetadata(total int64, page, pageSize int) Metadata {
	if total == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(total) / float64(pageSize))),
		TotalRecords: total,
	}
}
