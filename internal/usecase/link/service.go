package link

import (
	"context"
	"fmt"
	"strings"
	linkdomain "url-shortener/internal/domain"
	"url-shortener/internal/lib/random"
)

const (
	hashLength = 10
	maxRetries = 3
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, url string) (*linkdomain.Link, error) {
	url = strings.TrimSpace(url)
	if url == "" {
		return nil, ErrInvalidURL
	}

	var lastErr error

	for range maxRetries {
		hash := random.NewRandomString(hashLength)

		link := &linkdomain.Link{
			Url:  url,
			Hash: hash,
		}

		created, err := s.repo.Create(ctx, link)
		if err != nil {
			if strings.Contains(err.Error(), "hash collision") {
				lastErr = err
				continue
			}
			return nil, err
		}

		return created, nil
	}

	return nil, fmt.Errorf("failed to generate unique hash after %d attempts: %w", maxRetries, lastErr)
}

func (s *Service) GetByHash(ctx context.Context, hash string) (*linkdomain.Link, error) {
	hash = strings.TrimSpace(hash)

	if hash == "" {
		return nil, ErrInvalidHash
	}

	return s.repo.GetByHash(ctx, hash)

}
