package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/aeilang/urlshortener/config"
	"github.com/aeilang/urlshortener/internal/cache"
	"github.com/aeilang/urlshortener/internal/model"
	"github.com/aeilang/urlshortener/internal/repo"
	"github.com/aeilang/urlshortener/pkg/shortcode"
)

type URLService interface {
	CreateURL(ctx context.Context, params model.CreateURLRequest) (*repo.Url, error)
	GetURL(ctx context.Context, shortCode string) (*repo.Url, error)
	Cleanup(ctx context.Context) error
}

type urlService struct {
	querier   repo.Querier
	cache     cache.Cache
	generator shortcode.ShortCodeGenerator
	db        *sql.DB
	cfg       *config.Config
}

func NewURLService(db *sql.DB, cache cache.Cache, generator shortcode.ShortCodeGenerator, cfg *config.Config) URLService {
	s := urlService{
		querier:   repo.New(db),
		cache:     cache,
		generator: generator,
		db:        db,
		cfg:       cfg,
	}
	return &s
}

func (s *urlService) CreateURL(ctx context.Context, params model.CreateURLRequest) (*repo.Url, error) {
	var shortCode string
	var isCustom bool
	var expiresAt time.Time
	var err error

	if params.CustomCode != "" {
		// Check if custom code is available
		isAvailable, err := s.querier.IsShortCodeAvailable(ctx, params.CustomCode)
		if err != nil {
			return nil, err
		}
		if !isAvailable {
			return nil, errors.New("custom code is already taken")
		}
		shortCode = params.CustomCode
		isCustom = true
	} else {
		shortCode, err = s.tryFiveIsAvaliable(ctx, 0)
		if err != nil {
			return nil, err
		}
	}

	if params.Duration == nil {
		expiresAt = time.Now().Add(s.cfg.App.DefaultExpiration)
	} else {
		expiresAt = time.Now().Add(time.Hour * time.Duration(*params.Duration))
	}

	url, err := s.querier.CreateURL(ctx, repo.CreateURLParams{
		OriginalUrl: params.OriginalURL,
		ShortCode:   shortCode,
		ExpiresAt:   expiresAt,
		IsCustom:    isCustom,
	})

	if err != nil {
		return nil, err
	}

	// Cache the URL
	if err := s.cache.SetURL(ctx, url); err != nil {
		return nil, err
	}

	return &url, nil
}

func (s *urlService) tryFiveIsAvaliable(ctx context.Context, n int) (string, error) {
	if n >= 5 {
		return "", errors.New("try 5 times and failed")
	}
	shortCode := s.generator.NextID()

	isAvailale, err := s.querier.IsShortCodeAvailable(ctx, shortCode)
	if err != nil {
		return "", err
	}
	if !isAvailale {
		return s.tryFiveIsAvaliable(ctx, n+1)
	}
	return shortCode, nil
}

func (s *urlService) GetURL(ctx context.Context, shortCode string) (*repo.Url, error) {
	// Try cache first
	url, err := s.cache.GetURL(ctx, shortCode)
	if err != nil {
		return nil, err
	}
	if url != nil {
		return url, nil
	}

	// If not in cache, get from database
	url2, err := s.querier.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	// Cache the result
	if err := s.cache.SetURL(ctx, url2); err != nil {
		return nil, err
	}

	return url, nil
}

func (s *urlService) Cleanup(ctx context.Context) error {
	return s.querier.DeleteExpiredURLs(ctx)
}
