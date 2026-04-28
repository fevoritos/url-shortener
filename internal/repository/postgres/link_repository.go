package postgres

import (
	"context"
	"errors"
	"fmt"
	linkdomain "url-shortener/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, link *linkdomain.Link) (*linkdomain.Link, error) {
	const op = "repository.postgres.Create"

	query := `
		INSERT INTO links (url, hash)
		VALUES ($1, $2)
		ON CONFLICT (url) DO UPDATE SET url = EXCLUDED.url
		RETURNING url, hash, created_at, updated_at
	`

	res := &linkdomain.Link{}
	err := r.pool.QueryRow(ctx, query, link.Url, link.Hash).Scan(
		&res.Url,
		&res.Hash,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "links_hash_key" {
				return nil, fmt.Errorf("%s: %w: %w", op, linkdomain.ErrConflict, err)
			}
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

func (r *Repository) GetByHash(ctx context.Context, hash string) (*linkdomain.Link, error) {
	const op = "repository.postrges.GetByHash"

	query := `
		SELECT url, hash, created_at, updated_at FROM links
		WHERE hash = $1
	`

	res := &linkdomain.Link{}
	err := r.pool.QueryRow(ctx, query, hash).Scan(
		&res.Url,
		&res.Hash,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, linkdomain.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}
