package service

import (
	"context"
	"errors"
	"log/slog"
	"time"
	"url-shortener/internal/domain"
	"url-shortener/internal/repository"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
)

type LinkService struct {
	repo        *repository.LinkRepository
	redisClient *redis.Client
}

func NewLinkService(repo *repository.LinkRepository, redis *redis.Client) *LinkService {
	return &LinkService{
		repo:        repo,
		redisClient: redis,
	}
}

func (s *LinkService) CreateShortLink(ctx context.Context, longURL string) (domain.Link, error) {

	link := domain.Link{
		LongURL: longURL,
	}

	err := link.Validate()
	if err != nil {
		return domain.Link{}, err
	}

	//shortLink := domain.GenerateShortCode()

	//link.ShortCode = shortLink // положили в нашу струтуру

	//err = s.repo.Save(ctx, link) // сохраняем нашу  структуру в базе
	//if err != nil {
	//	return domain.Link{}, err
	//}

	for attempt := 0; attempt < 3; attempt++ {

		shortLink := domain.GenerateShortCode()

		link.ShortCode = shortLink

		err = s.repo.Save(ctx, link)
		if err == nil {
			return link, nil
		}

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "links_short_code_key" {
				continue
			}
		}

		return domain.Link{}, err
	}

	return domain.Link{}, errors.New("failed to generate unique short code after 3 attempts")

}

func (s *LinkService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {

	//  Ищем URL в Redis
	cacheURL, err := s.redisClient.Get(
		ctx, shortCode,
	).Result()

	// Cache Hit — ссылка найдена
	if err == nil {
		err = s.repo.IncrementClicks(ctx, shortCode)
		if err != nil {
			return "", err
		}

		return cacheURL, nil
	}

	// Redis недоступен — продолжаем через PostgreSQL
	if err != redis.Nil {
		slog.WarnContext(
			ctx,
			"Redis GET failed",
			"error", err,
		)
	}

	// увеличиваем счетчик и  получаем ссылку
	longURL, err := s.repo.IncrementClicksAndGetLink(ctx, shortCode)
	if err != nil {
		return "", err
	}

	//  Сохраняем URL в Redis на один час
	err = s.redisClient.Set(
		ctx,
		shortCode,
		longURL,
		time.Hour,
	).Err()

	if err != nil {
		slog.WarnContext(
			ctx,
			"Redis SET failed",
			"error", err,
		)
	}

	//Возвращаем оригинальный URL
	return longURL, nil
}

func (s *LinkService) GetStats(ctx context.Context, shortCode string) (domain.Link, error) {

	link, err := s.repo.FinByCode(ctx, shortCode)
	if err != nil {
		return domain.Link{}, err
	}

	return link, nil
}
