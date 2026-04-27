package link

import (
	"context"
	linkdomain "url-shortener/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, link *linkdomain.Link) (*linkdomain.Link, error)
	GetByHash(ctx context.Context, hash string) (*linkdomain.Link, error)
}

type Usecase interface {
	Create(ctx context.Context, url string) (*linkdomain.Link, error)
	GetByHash(ctx context.Context, hash string) (*linkdomain.Link, error)
}
