package auth

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthModule struct {
	handler *Handler
}

func NewAuthModule(pool *pgxpool.Pool) *AuthModule {
	repo := NewRepository(pool)
	tkManager := NewTokenManager([]byte("secret"), time.Minute*3)
	serv := NewService(repo, tkManager)
	hand := NewHandler(serv, tkManager)

	return &AuthModule{hand}
}

func (a *AuthModule) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", a.handler.Register)
	mux.HandleFunc("POST /auth/login", a.handler.Login)
	mux.HandleFunc("GET /auth/me", a.handler.Me)
}
