package repository

import (
	"context"
	"url-shortener/internal/database"
	"url-shortener/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LinkRepository struct {
	pool    *pgxpool.Pool
	queries *database.Queries
}

func NewLinkRepository(pool *pgxpool.Pool) *LinkRepository {

	return &LinkRepository{pool: pool,
		queries: database.New(pool)}
}

func (l *LinkRepository) Save(ctx context.Context, link domain.Link) error {

	params := database.CreateLinkParams{
		ShortCode: link.ShortCode,
		LongUrl:   link.LongURL,
	}

	_, err := l.queries.CreateLink(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (l *LinkRepository) FinByCode(ctx context.Context, shortCode string) (domain.Link, error) {

	dbLink, err := l.queries.GetLinkByCode(ctx, shortCode)
	if err != nil {
		return domain.Link{}, err
	}

	link := domain.Link{
		ID:        int(dbLink.ID),
		ShortCode: dbLink.ShortCode,
		LongURL:   dbLink.LongUrl,
		Clicks:    int(dbLink.Clicks),
		CreatedAt: dbLink.CreatedAt.Time,
	}

	return link, nil
}

func (l *LinkRepository) IncrementClicks(ctx context.Context, shortCode string) error {

	err := l.queries.IncrementClicks(ctx, shortCode)
	if err != nil {
		return err
	}

	return nil
}

func (l *LinkRepository) IncrementClicksAndGetLink(ctx context.Context, shortCode string) (string, error) {

	longURL, err := l.queries.IncrementClicksAndGetLink(
		ctx,
		shortCode,
	)
	if err != nil {
		return "", err
	}

	return longURL, nil
}
