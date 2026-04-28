package link

import "errors"

var (
	ErrConflict = errors.New("data conflict")
	ErrNotFound = errors.New("not found")
)
