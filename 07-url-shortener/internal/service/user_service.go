package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strconv"
	"time"
	"url-shortener/internal/database"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidRegistrationData = errors.New("invalid registration data")
var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrInvalidCredentials = errors.New("invalid email or password")

type UserStore interface {
	CreateUser(ctx context.Context, arg database.CreateUserParams) (database.CreateUserRow, error)
	GetUserByEmail(ctx context.Context, email string) (database.User, error)
}

type UserService struct {
	repo      UserStore
	jwtSecret string
}

func NewUserService(repo UserStore, jwtSecret string) *UserService {
	return &UserService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (u *UserService) Register(ctx context.Context, email string, password string) error {

	if email == "" {
		return fmt.Errorf("%w: email is required", ErrInvalidRegistrationData)
	}

	if password == "" {
		return fmt.Errorf("%w: password is required", ErrInvalidRegistrationData)
	}

	if len(password) < 8 {
		return fmt.Errorf(
			"%w: password must be at least 8 characters", ErrInvalidRegistrationData)
	}

	address, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("%w: invalid email format", ErrInvalidRegistrationData)
	}

	if address.Address != email {
		return fmt.Errorf("%w: invalid email format", ErrInvalidRegistrationData)
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	params := database.CreateUserParams{
		Email:        email,
		PasswordHash: string(hash),
	}

	_, err = u.repo.CreateUser(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" && // нарушение уникальности поля
			pgErr.ConstraintName == "users_email_key" { // ограничение уникальности почты

			return ErrEmailAlreadyExists
		}

		return fmt.Errorf(
			"failed to create user: %w", err,
		)
	}

	return nil
}

func (u *UserService) Login(ctx context.Context, email string, password string) (string, error) {

	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword( // проверка пороля
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("failed to verify password: %w", err)
	}

	claims := jwt.RegisteredClaims{
		Subject: strconv.FormatInt(user.ID, 10),
		ExpiresAt: jwt.NewNumericDate( // преобразуем в формат который ожидает наша библиотека   и сохраняем время окончания токена
			time.Now().Add(1 * time.Hour), // текущее время  + добавляем к нему час
		),
	}

	token := jwt.NewWithClaims( // создаем объект токена , с нашими данными
		jwt.SigningMethodHS256, // алгоритм подписи
		claims,                 // время окончания действия токена и айди
	)

	tokenString, err := token.SignedString( // подписывает  JWT нашим сикретным ключом
		[]byte(u.jwtSecret),
	)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return tokenString, nil
}
