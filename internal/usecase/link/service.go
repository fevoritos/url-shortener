package link

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
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
	if url == "" || !isValidURL(url) {
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
			if errors.Is(err, linkdomain.ErrConflict) {
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

func isValidURL(rawURL string) bool {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	host, port, err := net.SplitHostPort(u.Host)
	hasPort := err == nil
	if !hasPort {
		host = u.Host
	}

	if host == "" {
		return false
	}

	if hasPort && port != "" {
		return true
	}

	dotIdx := strings.LastIndex(host, ".")
	if dotIdx == -1 || dotIdx == 0 || dotIdx == len(host)-1 {
		return false
	}

	return !strings.Contains(host, " ")
}
