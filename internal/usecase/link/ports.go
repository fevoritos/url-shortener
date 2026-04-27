package link

import (
	"context"
	linkdomain "url-shortener/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, Url string) (*linkdomain.Link, error)
	GetByHash(ctx context.Context, Hash string) (*linkdomain.Link, error)
}

type Usecase interface {
	Create(ctx context.Context, Url string) (*linkdomain.Link, error)
	GetByHash(ctx context.Context, Hash string) (*linkdomain.Link, error)
}
