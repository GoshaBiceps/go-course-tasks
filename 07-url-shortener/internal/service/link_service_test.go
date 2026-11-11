package service

import (
	"context"
	"testing"
)

func TestCreateShortLink_EmptyURL(t *testing.T) {
	svc := &LinkService{}

	_, err := svc.CreateShortLink(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}

func TestCreateShortLink_InvalidScheme(t *testing.T) {
	svc := &LinkService{}

	_, err := svc.CreateShortLink(context.Background(), "ftp://google.com")
	if err == nil {
		t.Fatal("expected error for invalid URL scheme")
	}
}
