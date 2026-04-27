package postgres

import (
	"context"
	linkdomain "url-shortener/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, Url string) (*linkdomain.Link, error) {
	return &linkdomain.Link{}, nil
}
func (r *Repository) GetByHash(ctx context.Context, Hash string) (*linkdomain.Link, error) {
	return &linkdomain.Link{}, nil
}
