package user

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type UserHandler interface {
	GetLatestCommit(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	user   UserRepository
	logger *zap.SugaredLogger
}

func NewHandler(ctx context.Context) *Handler {
	logger, ok := ctx.Value("logger").(*zap.SugaredLogger)
	if !ok {
		log.Fatalf("Could not retrieve logger from context")
	}
	postgres, ok := ctx.Value("postgres").(*sqlx.DB)
	if !ok {
		log.Fatalf("Could not retrieve router from context")
	}
	user := NewRepository(postgres)
	return &Handler{user, logger}
}

func (h *Handler) Routes(r *chi.Mux) chi.Router {
	r.Get("/latest-commit", h.GetLatestCommit)
	return r
}

func (h *Handler) GetLatestCommit(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
