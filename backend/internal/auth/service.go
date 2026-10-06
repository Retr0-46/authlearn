package auth

import (
	"context"
	"strings"
)

type ServiceInterface interface {
	Register(ctx context.Context, username, password string) (User, error)
	Login(ctx context.Context, username, password string) (LoginResult, error)
	GetUserInfo(ctx context.Context, userID int64) (User, error)
}

type Service struct {
	repo         RepoInterface
	tokenManager TokenManagerInterface
}

func NewService(repository RepoInterface, tokenManager TokenManagerInterface) *Service {
	return &Service{repo: repository, tokenManager: tokenManager}
}

func (s *Service) Register(ctx context.Context, username, password string) (User, error) {

	name := strings.TrimSpace(username)

	if len(name) < 1 {
		return User{}, ErrEmptyUsername
	}

	if len([]rune(password)) < 4 {
		return User{}, ErrTooShortPassword
	}

	if len([]byte(password)) > 72 {
		return User{}, ErrTooLongPassword
	}

	hashPassword, err := HashPassword(password)
	if err != nil {
		return User{}, err
	}

	user, err := s.repo.CreateUser(ctx, name, hashPassword)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, username, password string) (LoginResult, error) {
	user, err := s.repo.GetUser(ctx, username)

	if err != nil {
		return LoginResult{}, err
	}

	if ComparePassword(user.PasswordHash, password) != nil {
		return LoginResult{}, ErrPasswordMismatch
	}

	token, err := s.tokenManager.IssueAccessToken(user)

	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{user, token.Value, token.ExpiresAt}, err
}

func (s *Service) GetUserInfo(ctx context.Context, userID int64) (User, error) {
	user, err := s.repo.GetUser(ctx, userID)

	if err != nil {
		return User{}, err
	}

	return user, nil
}
