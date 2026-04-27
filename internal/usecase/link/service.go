package link

import (
	"context"
	"fmt"
	"strings"
	linkdomain "url-shortener/internal/domain"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, Url string) (*linkdomain.Link, error) {
	Url = strings.TrimSpace(Url)

	if Url == "" {
		return &linkdomain.Link{}, fmt.Errorf("%w: Url is required", ErrInvalidURL)
	}

	created, err := s.repo.Create(ctx, Url)
	if err != nil {
		return nil, err
	}

	return created, nil

}

func (s *Service) GetByHash(ctx context.Context, Hash string) (*linkdomain.Link, error) {
	Hash = strings.TrimSpace(Hash)

	if Hash == "" {
		return &linkdomain.Link{}, fmt.Errorf("%w: ", ErrInvalidHash)
	}

	return s.repo.GetByHash(ctx, Hash)

}
