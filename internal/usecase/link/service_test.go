package link

import (
	"context"
	"testing"

	linkdomain "url-shortener/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, link *linkdomain.Link) (*linkdomain.Link, error) {
	args := m.Called(ctx, link)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*linkdomain.Link), args.Error(1)
}

func (m *MockRepository) GetByHash(ctx context.Context, hash string) (*linkdomain.Link, error) {
	args := m.Called(ctx, hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*linkdomain.Link), args.Error(1)
}

func TestService_Create(t *testing.T) {
	ctx := context.Background()
	rawURL := "https://github.com"

	t.Run("success_create", func(t *testing.T) {
		repo := new(MockRepository)
		service := NewService(repo)
		rawURL := "https://google.com"

		repo.On("Create", mock.Anything, mock.MatchedBy(func(link *linkdomain.Link) bool {
			return link.Url == rawURL && len(link.Hash) == 10
		})).Return(&linkdomain.Link{
			Url:  rawURL,
			Hash: "randomhash",
		}, nil)

		link, err := service.Create(ctx, rawURL)

		assert.NoError(t, err)
		assert.NotNil(t, link)
		assert.Equal(t, rawURL, link.Url)
		assert.Len(t, link.Hash, 10)

		repo.AssertExpectations(t)
	})

	t.Run("invalid_url", func(t *testing.T) {
		repo := new(MockRepository)
		service := NewService(repo)

		invalidURLs := []string{"322", "https://google.", "not-a-url", ""}

		for _, u := range invalidURLs {
			link, err := service.Create(ctx, u)
			assert.ErrorIs(t, err, ErrInvalidURL)
			assert.Nil(t, link)
		}
	})

	t.Run("hash_collision_retry_success", func(t *testing.T) {
		repo := new(MockRepository)
		service := NewService(repo)

		repo.On("Create", mock.Anything, mock.MatchedBy(func(l *linkdomain.Link) bool {
			return l.Url == rawURL
		})).Return(nil, linkdomain.ErrConflict).Once()

		repo.On("Create", mock.Anything, mock.MatchedBy(func(l *linkdomain.Link) bool {
			return l.Url == rawURL
		})).Return(&linkdomain.Link{
			Url:  rawURL,
			Hash: "newhash123",
		}, nil).Once()

		link, err := service.Create(ctx, rawURL)

		assert.NoError(t, err)
		assert.NotNil(t, link)
		assert.Equal(t, rawURL, link.Url)

		repo.AssertExpectations(t)
	})
}

func TestService_GetByHash(t *testing.T) {
	ctx := context.Background()

	t.Run("found", func(t *testing.T) {
		repo := new(MockRepository)
		service := NewService(repo)
		hash := "ABC123456_"

		repo.On("GetByHash", ctx, hash).Return(&linkdomain.Link{Url: "https://ya.ru", Hash: hash}, nil)

		link, err := service.GetByHash(ctx, hash)

		assert.NoError(t, err)
		assert.Equal(t, "https://ya.ru", link.Url)
	})

	t.Run("not_found", func(t *testing.T) {
		repo := new(MockRepository)
		service := NewService(repo)

		repo.On("GetByHash", ctx, "missing").Return(nil, linkdomain.ErrNotFound)

		link, err := service.GetByHash(ctx, "missing")

		assert.ErrorIs(t, err, linkdomain.ErrNotFound)
		assert.Nil(t, link)
	})
}

func Test_isValidURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{"valid https", "https://google.com", true},
		{"valid http", "http://my-site.ru", true},
		{"valid with port", "http://localhost:8080", true},
		{"invalid protocol", "ftp://files.com", false},
		{"no protocol", "google.com", false},
		{"just string", "322", false},
		{"trailing dot", "https://google.", false},
		{"no dot in host", "https://localhost", false},
		{"with space", "https://goo gle.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidURL(tt.url); got != tt.want {
				t.Errorf("isValidURL() = %v, want %v for %s", got, tt.want, tt.url)
			}
		})
	}
}
