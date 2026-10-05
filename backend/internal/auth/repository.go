package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepoInterface interface {
	CreateUser(ctx context.Context, username string, passwordHash string) (User, error)
	GetUser(ctx context.Context, identifier any) (User, error)
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(p *pgxpool.Pool) *Repository {
	return &Repository{pool: p}
}

func (r *Repository) CreateUser(ctx context.Context, username string, passwordHash string) (User, error) {
	var user User

	row := r.pool.QueryRow(ctx, `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, username, password_hash, is_active, created_at, updated_at`, username, passwordHash)
	err := row.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_username_unique" {
		return User{}, ErrUserAlreadyExist
	}

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *Repository) GetUser(ctx context.Context, identifier any) (User, error) {
	switch t := identifier.(type) {
	case string:
		return r.getUserByUsername(ctx, t)
	case int:
		return r.getUserByID(ctx, int64(t))
	case int64:
		return r.getUserByID(ctx, t)
	default:
		return User{}, ErrWrongIdentificator
	}
}

func (r *Repository) getUserByID(ctx context.Context, id int64) (User, error) {
	var user User

	row := r.pool.QueryRow(ctx, `SELECT id, username, password_hash, is_active, created_at, updated_at FROM users WHERE id = $1`, id)
	err := row.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *Repository) getUserByUsername(ctx context.Context, name string) (User, error) {
	var user User

	row := r.pool.QueryRow(ctx, `SELECT id, username, password_hash, is_active, created_at, updated_at FROM users WHERE username = $1`, name)
	err := row.Scan(&user.Id, &user.Username, &user.PasswordHash, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}

	if err != nil {
		return User{}, err
	}

	return user, nil
}


