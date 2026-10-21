package service

import (
	"context"
	"log/slog"
	"time"
	"url-shortener/internal/domain"
	"url-shortener/internal/repository"

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

	shortLink := domain.GenerateShortCode()

	link.ShortCode = shortLink // положили в нашу струтуру

	err = s.repo.Save(ctx, link) // сохраняем нашу  структуру в базе
	if err != nil {
		return domain.Link{}, err
	}

	return link, nil
}

func (s *LinkService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {

	// 1. Ищем URL в Redis
	cacheURL, err := s.redisClient.Get(
		ctx, shortCode,
	).Result()

	// 2. Cache Hit — ссылка найдена
	if err == nil {
		err = s.repo.IncrementClicks(ctx, shortCode)
		if err != nil {
			return "", err
		}

		return cacheURL, nil
	}

	// 3. Redis недоступен — продолжаем через PostgreSQL
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

	// 6. Сохраняем URL в Redis на один час
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

	// 7. Возвращаем оригинальный URL
	return longURL, nil
}

func (s *LinkService) GetStats(ctx context.Context, shortCode string) (domain.Link, error) {

	link, err := s.repo.FinByCode(ctx, shortCode)
	if err != nil {
		return domain.Link{}, err
	}

	return link, nil
}
