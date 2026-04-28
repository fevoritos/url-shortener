package memory

import (
	"context"
	"errors"
	"sync"

	linkdomain "url-shortener/internal/domain"
)

type Repository struct {
	mu        sync.RWMutex
	hashToUrl map[string]*linkdomain.Link
	urlToHash map[string]*linkdomain.Link
}

func New() *Repository {
	return &Repository{
		hashToUrl: make(map[string]*linkdomain.Link),
		urlToHash: make(map[string]*linkdomain.Link),
	}
}

func (r *Repository) Create(ctx context.Context, link *linkdomain.Link) (*linkdomain.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if existing, ok := r.urlToHash[link.Url]; ok {
		return existing, nil
	}

	if _, ok := r.hashToUrl[link.Hash]; ok {
		return nil, errors.New("hash collision")
	}

	r.hashToUrl[link.Hash] = link
	r.urlToHash[link.Url] = link

	return link, nil
}

func (r *Repository) GetByHash(ctx context.Context, hash string) (*linkdomain.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	link, ok := r.hashToUrl[hash]
	if !ok {
		return nil, errors.New("not found")
	}

	return link, nil
}
