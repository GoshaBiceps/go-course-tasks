package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"url-shortener/internal/database"
	"url-shortener/internal/repository"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

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
		return errors.New("email is required")
	}

	if password == "" {
		return errors.New("password is required")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	address, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email format")
	}

	if address.Address != email {
		return errors.New("invalid email format")
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
		return fmt.Errorf("failed to create user: %w", err)
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
