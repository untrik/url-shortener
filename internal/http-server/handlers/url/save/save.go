package save

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"

	resp "github.com/untrik/url-shortener/internal/lib/api/response"
	sh "github.com/untrik/url-shortener/internal/lib/api/save-helper"
	"github.com/untrik/url-shortener/internal/lib/logger/sl"
	"github.com/untrik/url-shortener/internal/lib/random"
	"github.com/untrik/url-shortener/storage"
)

const aliasLength = 8

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode Request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode Request"))
			return
		}
		log.Info("Request body decoded", slog.Any("request", req))
		if err := validator.New().Struct(req); err != nil {
			validateError := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(validateError))
			return
		}
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}
		id, err := urlSaver.SaveURL(req.URL, alias)
		if aliasUniqueErr, ok := err.(*pq.Error); ok && aliasUniqueErr.Code == "23505" {
			id, alias, err = sh.SaveWithRetry(urlSaver, req.URL, aliasLength)
			if err != nil {
				log.Error("alias save error", sl.Err(err))
				render.JSON(w, r, resp.Error("failed to add url"))
				return
			}
		}
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("URL already exists", slog.String("url", req.URL))
			render.JSON(w, r, resp.Error("URL already exists"))
			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}
		log.Info("url added", slog.Int64("id", id))
		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
