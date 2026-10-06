package repository

import (
	"context"
	"errors"
	"url-shortener/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type URLRepository struct {
	db *pgxpool.Pool
}

func NewURLRepository(db *pgxpool.Pool) *URLRepository {
	return &URLRepository{
		db: db,
	}
}

func (r *URLRepository) Create(ctx context.Context, shortCode string, originalURL string) error {

	_, err := r.db.Exec(
		ctx,
		"INSERT INTO urls (short_code, original_url) VALUES ($1, $2)",
		shortCode,
		originalURL,
	)
	return err
}

func (r *URLRepository) GetByShortCode(ctx context.Context, shortCode string) (*models.URL, error) {
	var url models.URL
	err := r.db.QueryRow(
		ctx,
		"SELECT id, short_code, original_url, access_count + 1, created_at, updated_at FROM urls WHERE short_code = $1", shortCode,
	).Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.AccessCount, &url.CreatedAt, &url.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, pgx.ErrNoRows
	}

	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *URLRepository) Update(ctx context.Context, shortCode string, newURL string) (int64, error) {
	result, err := r.db.Exec(
		ctx,
		"UPDATE urls SET original_url = $1, access_count = access_count + 1, updated_at = NOW() WHERE short_code = $2;", newURL, shortCode)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *URLRepository) Delete(ctx context.Context, shortCode string) (int64, error) {
	result, err := r.db.Exec(
		ctx,
		"DELETE FROM urls WHERE short_code = $1", shortCode)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *URLRepository) IncrementAccessCount(
	ctx context.Context,
	shortCode string,
) error {

	_, err := r.db.Exec(
		ctx,
		`UPDATE urls
		 SET access_count = access_count + 1
		 WHERE short_code = $1`,
		shortCode,
	)

	return err
}
