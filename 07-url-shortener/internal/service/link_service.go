package service

import (
	"context"
	"url-shortener/internal/domain"
	"url-shortener/internal/repository"
)

type LinkService struct {
	repo *repository.LinkRepository
}

func NewLinkService(repo *repository.LinkRepository) *LinkService {
	return &LinkService{
		repo: repo,
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

	link, err := s.repo.FinByCode(ctx, shortCode)
	if err != nil {
		return "", err
	}

	err = s.repo.IncrementClicks(ctx, shortCode)
	if err != nil {
		return "", err
	}

	return link.LongURL, nil
}

func (s *LinkService) GetStats(ctx context.Context, shortCode string) (domain.Link, error) {

	link, err := s.repo.FinByCode(ctx, shortCode)
	if err != nil {
		return domain.Link{}, err
	}

	return link, nil
}
