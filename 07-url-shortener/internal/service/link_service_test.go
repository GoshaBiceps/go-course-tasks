package service

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"
	"url-shortener/internal/domain"

	"github.com/alicebob/miniredis/v2"
	"github.com/alicebob/miniredis/v2/server"
	"github.com/redis/go-redis/v9"
)

type MockLinkStore struct {
	saveErr         error
	saveCalled      bool        // Вызывался ли Save
	savedLink       domain.Link // Ссылка, переданная сервисом на сохранение
	foundLink       domain.Link // Ссылка, которую мок возвращает вместо PostgreSQL
	findErr         error       // Ошибка поиска ссылки
	findCalled      bool        // Вызывался ли FinByCode
	findCode        string      // Код, по которому сервис искал ссылку
	incrementCalled bool        // Вызывался ли IncrementClicks
	incrementErr    error       // Ошибка, возвращаемая методом IncrementClicks
	getLinkCalled   bool        // Вызывался ли IncrementClicksAndGetLink
	getLinkCode     string      // Код, по которому сервис получал исходный URL
	getLinkErr      error       // Ошибка, возвращаемая методом IncrementClicksAndGetLink
}

func (m *MockLinkStore) Save(ctx context.Context, link domain.Link) error {
	m.saveCalled = true
	m.savedLink = link
	return m.saveErr
}

func (m *MockLinkStore) FinByCode(ctx context.Context, shortCode string) (domain.Link, error) {
	m.findCalled = true
	m.findCode = shortCode
	return m.foundLink, m.findErr
}

func (m *MockLinkStore) IncrementClicks(ctx context.Context, shortCode string) error {
	m.incrementCalled = true
	return m.incrementErr
}

func (m *MockLinkStore) IncrementClicksAndGetLink(ctx context.Context, shortCode string) (string, error) {
	m.getLinkCalled = true
	m.getLinkCode = shortCode
	return m.foundLink.LongURL, m.getLinkErr
}

func TestCreateShortLink_EmptyURL(t *testing.T) {
	mock := &MockLinkStore{}
	svc := NewLinkService(mock, nil)

	_, err := svc.CreateShortLink(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}

	if mock.saveCalled {
		t.Fatal("Save must not be called for empty URL")
	}
}

func TestCreateShortLink_InvalidScheme(t *testing.T) {
	mock := &MockLinkStore{}
	svc := NewLinkService(mock, nil)

	_, err := svc.CreateShortLink(context.Background(), "ftp://google.com")
	if err == nil {
		t.Fatal("expected error for invalid URL scheme")
	}

	if mock.saveCalled {
		t.Fatal("Save must not be called for invalid URL scheme")
	}
}

func TestCreateShortLink_Success(t *testing.T) {
	mock := &MockLinkStore{}
	svc := NewLinkService(mock, nil)

	originalURL := "https://google.com"
	link, err := svc.CreateShortLink(context.Background(), originalURL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if link.ShortCode == "" {
		t.Fatal("short code is empty")
	}

	if link.LongURL != originalURL {
		t.Fatalf("expected %s, got %s", originalURL, link.LongURL)
	}

	if !mock.saveCalled {
		t.Fatal("expected Save to be called")
	}

	// Проверяем, что репозиторий получил правильные данные для сохранения.
	if mock.savedLink.LongURL != originalURL {
		t.Fatalf("expected saved URL %s, got %s", originalURL, mock.savedLink.LongURL)
	}

	if mock.savedLink.ShortCode != link.ShortCode {
		t.Fatalf("expected saved short code %s, got %s", link.ShortCode, mock.savedLink.ShortCode)
	}
}

func TestCreateShortLink_SaveError(t *testing.T) {
	expectedErr := errors.New("database error")
	mock := &MockLinkStore{
		saveErr: expectedErr,
	}
	svc := NewLinkService(mock, nil)

	_, err := svc.CreateShortLink(context.Background(), "https://google.com")
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if !mock.saveCalled {
		t.Fatal("expected Save to be called")
	}
}

func TestGetStats_Success(t *testing.T) {
	mock := &MockLinkStore{
		foundLink: domain.Link{
			ShortCode: "abc123",
			LongURL:   "https://google.com",
			Clicks:    15,
		},
	}
	svc := NewLinkService(mock, nil)

	link, err := svc.GetStats(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if link.ShortCode != "abc123" {
		t.Fatalf("expected abc123, got %s", link.ShortCode)
	}

	if link.Clicks != 15 {
		t.Fatalf("expected 15 clicks, got %v", link.Clicks)
	}

	if link.LongURL != "https://google.com" {
		t.Fatalf("expected https://google.com, got %s", link.LongURL)
	}

	if !mock.findCalled {
		t.Fatal("expected FinByCode to be called")
	}

	if mock.findCode != "abc123" {
		t.Fatalf("expected code abc123, got %s", mock.findCode)
	}
}

func TestGetStats_Error(t *testing.T) {
	expectedErr := errors.New("link not found")
	mock := &MockLinkStore{
		findErr: expectedErr,
	}
	svc := NewLinkService(mock, nil)

	_, err := svc.GetStats(context.Background(), "abc123")
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if !mock.findCalled {
		t.Fatal("expected FinByCode to be called")
	}
}

func TestGetOriginalURL_CacheHit(t *testing.T) {
	// mr принимает команды Redis и хранит тестовые данные.
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	if err := mr.Set("abc123", "https://google.com"); err != nil {
		t.Fatalf("failed to seed miniredis: %v", err)
	}

	// Клиент отправляет команды GET и SET на адрес тестового Redis.
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer client.Close()

	mock := &MockLinkStore{}
	svc := NewLinkService(mock, client) // Создали сервис с тестовыми зависимостями.

	url, err := svc.GetOriginalURL(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if url != "https://google.com" {
		t.Fatalf("expected https://google.com, got %s", url)
	}

	if !mock.incrementCalled { // Проверяем, вызывался ли IncrementClicks.
		t.Fatal("expected IncrementClicks to be called")
	}

	if mock.getLinkCalled {
		t.Fatal("IncrementClicksAndGetLink must not be called on cache hit")
	}
}

func TestGetOriginalURL_CacheMiss(t *testing.T) {
	mr, err := miniredis.Run() // Запускаем тестовый Redis; кеш пока пустой.
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer client.Close()

	mock := &MockLinkStore{
		foundLink: domain.Link{
			LongURL: "https://google.com",
		},
	}
	svc := NewLinkService(mock, client)

	url, err := svc.GetOriginalURL(context.Background(), "abc123") // Вызываем настоящий метод сервиса.
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if url != "https://google.com" { // Проверяем возвращённый URL.
		t.Fatalf("expected https://google.com, got %s", url)
	}

	if !mock.getLinkCalled { // Проверяем вызов метода репозитория.
		t.Fatal("expected IncrementClicksAndGetLink to be called")
	}

	if mock.getLinkCode != "abc123" {
		t.Fatalf("expected code abc123, got %s", mock.getLinkCode)
	}

	if mock.incrementCalled {
		t.Fatal("IncrementClicks must not be called on cache miss")
	}

	// Проверяем, что ссылка действительно сохранилась в кеш.
	cachedURL, err := client.Get(context.Background(), "abc123").Result()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cachedURL != "https://google.com" { // Проверяем сохранённое значение.
		t.Fatalf("expected https://google.com, got %s", cachedURL)
	}

	if ttl := mr.TTL("abc123"); ttl != time.Hour {
		t.Fatalf("expected cache TTL %v, got %v", time.Hour, ttl)
	}
}

func TestGetOriginalURL_IncrementClicksError(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	if err := mr.Set("abc123", "https://google.com"); err != nil {
		t.Fatalf("failed to seed miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer client.Close()

	expectedErr := errors.New("db error")
	mock := &MockLinkStore{
		incrementErr: expectedErr,
	}
	svc := NewLinkService(mock, client)

	url, err := svc.GetOriginalURL(context.Background(), "abc123")
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if url != "" {
		t.Fatalf("expected empty URL, got %s", url)
	}

	if !mock.incrementCalled {
		t.Fatal("expected IncrementClicks to be called")
	}

	if mock.getLinkCalled {
		t.Fatal("IncrementClicksAndGetLink must not be called on cache hit")
	}
}

func TestGetOriginalURL_DBError(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer client.Close()

	expectedErr := errors.New("db error")
	mock := &MockLinkStore{
		getLinkErr: expectedErr,
	}
	svc := NewLinkService(mock, client)

	url, err := svc.GetOriginalURL(context.Background(), "abc123")
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}

	if url != "" {
		t.Fatalf("expected empty URL, got %s", url)
	}

	if !mock.getLinkCalled {
		t.Fatal("expected IncrementClicksAndGetLink to be called")
	}

	if mr.Exists("abc123") {
		t.Fatal("URL must not be cached after database error")
	}
}

func TestGetOriginalURL_RedisUnavailable(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	mr.Close() // Отключаем Redis, но клиент уже знает его адрес.
	defer client.Close()

	mock := &MockLinkStore{
		foundLink: domain.Link{
			LongURL: "https://google.com",
		},
	}
	svc := NewLinkService(mock, client)

	url, err := svc.GetOriginalURL(context.Background(), "abc123")

	// Сервис не должен возвращать ошибку из-за недоступности Redis.
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Проверяем возвращённую ссылку.
	if url != "https://google.com" {
		t.Fatalf("expected https://google.com, got %s", url)
	}

	// Проверяем обращение к репозиторию.
	if !mock.getLinkCalled {
		t.Fatal("expected IncrementClicksAndGetLink to be called")
	}
}

func TestGetOriginalURL_RedisSetError(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer mr.Close()

	// Флажок: пытался ли сервис выполнить SET.
	var setCalled atomic.Bool

	// Перехватываем команды Redis: GET работает, а SET возвращает ошибку.
	mr.Server().SetPreHook(func(peer *server.Peer, cmd string, args ...string) bool {
		if strings.EqualFold(cmd, "SET") {
			setCalled.Store(true)
			peer.WriteError("ERR simulated SET failure")
			return true
		}
		return false
	})

	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer client.Close()

	mock := &MockLinkStore{
		foundLink: domain.Link{
			LongURL: "https://google.com",
		},
	}
	svc := NewLinkService(mock, client)

	url, err := svc.GetOriginalURL(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if url != "https://google.com" {
		t.Fatalf("expected https://google.com, got %s", url)
	}

	if !mock.getLinkCalled {
		t.Fatal("expected IncrementClicksAndGetLink to be called")
	}

	if !setCalled.Load() {
		t.Fatal("expected Redis SET to be called")
	}

	if mr.Exists("abc123") {
		t.Fatal("expected URL not to be cached")
	}
}
