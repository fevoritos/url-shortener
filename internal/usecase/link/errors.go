package link

import "errors"

var (
	ErrInvalidURL  = errors.New("invalid url")
	ErrInvalidHash = errors.New("invalid hash")
)
