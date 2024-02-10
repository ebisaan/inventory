package domain

import (
	"errors"
)

var (
	ErrNotFound            = errors.New("resource not found")
	ErrAssociationNotFound = errors.New("association resource not found")
	ErrEditConflict        = errors.New("edit conflicted")
)
