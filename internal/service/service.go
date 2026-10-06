package service

import (
	"context"

	"url-shortener/internal/models"
	"url-shortener/internal/repository"
	"url-shortener/internal/utils"

	"github.com/jackc/pgx/v5"
)

type URLService struct {
	repo *repository.URLRepository
}

func NewURLService(repo *repository.URLRepository) *URLService {
	return &URLService{
		repo: repo,
	}
}

func (s *URLService) CreateURL(
	ctx context.Context,
	originalURL string,
) (string, error) {

	shortCode := utils.GenerateShortCode(6)

	err := s.repo.Create(
		ctx,
		shortCode,
		originalURL,
	)

	if err != nil {
		return "", err
	}

	return shortCode, nil
}

func (s *URLService) GetURL(
	ctx context.Context,
	shortCode string,
) (*models.URL, error) {

	url, err := s.repo.GetByShortCode(ctx, shortCode)

	if err != nil {
		return nil, err
	}

	err = s.repo.IncrementAccessCount(ctx, shortCode)

	if err != nil {
		return nil, err
	}

	url.AccessCount++

	return url, nil
}

func (s *URLService) UpdateURL(
	ctx context.Context,
	shortCode string,
	originalURL string,
) error {

	rowsAffected, err := s.repo.Update(
		ctx,
		shortCode,
		originalURL,
	)

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (s *URLService) DeleteURL(
	ctx context.Context,
	shortCode string,
) error {

	rowsAffected, err := s.repo.Delete(
		ctx,
		shortCode,
	)

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
