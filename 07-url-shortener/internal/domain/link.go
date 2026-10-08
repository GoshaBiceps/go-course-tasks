package domain

import (
	"errors"
	"math/rand/v2"
	"net/url"
	"time"
)

type Link struct {
	ID        int
	ShortCode string
	LongURL   string
	Clicks    int
	CreatedAt time.Time
}

func GenerateShortCode() string {

	charest := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6) //  хранить код

	for i := range code {

		randomIndex := rand.IntN(len(charest)) // рандомный индекс
		code[i] = charest[randomIndex]
	}

	return string(code)
}

func (l Link) Validate() error {

	if l.LongURL == "" {
		return errors.New("lond URL is empty")
	}

	parsedURL, err := url.ParseRequestURI(l.LongURL)

	if err != nil {
		return errors.New("invalid URL")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("URL must use http or https")
	}

	return nil
}
