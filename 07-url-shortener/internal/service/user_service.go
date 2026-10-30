package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"url-shortener/internal/database"
	"url-shortener/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidRegistrationData = errors.New("invalid registration data")
var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrInvalidCredentials = errors.New("invalid email or password")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
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
			pgErr.Code == "23505" && // уникльность поля
			pgErr.ConstraintName == "users_email_key" { // ограничение уникальности почты

			return ErrEmailAlreadyExists
		}

		return fmt.Errorf(
			"failed to create user: %w", err,
		)
	}

	return nil
}

func (u *UserService) Login(ctx context.Context, email string, password string) error {

	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("invalid email or password")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errors.New("invalid email or password")
		}
		return fmt.Errorf("failed to verify password: %w", err)
	}
	return nil
}
