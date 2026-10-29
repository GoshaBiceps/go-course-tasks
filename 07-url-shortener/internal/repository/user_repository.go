package repository

import (
	"context"
	"url-shortener/internal/database"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	queries *database.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		queries: database.New(pool),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.CreateUserRow, error) {

	user, err := r.queries.CreateUser(ctx, arg)
	if err != nil {
		return database.CreateUserRow{}, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (database.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return database.User{}, err
	}
	return user, nil
}
