package delete

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/untrik/url-shortener/internal/lib/api/response"
	"github.com/untrik/url-shortener/internal/lib/logger/sl"
	"github.com/untrik/url-shortener/storage"
)

type DeleteURL interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, deleteUrl DeleteURL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, response.Error("invalid request"))
			return
		}
		err := deleteUrl.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("The URL for this alias was not found", "alias", alias)
			render.JSON(w, r, response.Error("The URL for this alias was not found"))
			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err))
			render.JSON(w, r, response.Error("internal error"))
			return
		}
		render.JSON(w, r, http.StatusOK)
	}
}
