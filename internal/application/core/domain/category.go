package domain

type MainCategory struct {
	ID   int64
	Name string
}

type SubCategory struct {
	ID           int64
	MainCategory string
	Name         string
}
