package delivery

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/s3nn1k/ef-mob-task/internal/models"
	"github.com/s3nn1k/ef-mob-task/internal/service"
	"github.com/s3nn1k/ef-mob-task/pkg/logger"
)

type Handler struct {
	log     *slog.Logger
	service *service.Service
}

func NewHandler(l *slog.Logger, s *service.Service) *Handler {
	return &Handler{
		log:     l,
		service: s,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var song models.Song

	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		h.response(w, Error("Can't decode json body"), http.StatusBadRequest)
		return
	}

	ctx := logger.NewCtxWithLog(r.Context(), h.log)

	song, err = h.service.Create(ctx, song.Song, song.Group)
	if err != nil {
		h.log.Error(err.Error(), "input", song.AsLogValue())

		h.response(w, Error("Can't create song"), http.StatusInternalServerError)
		return
	}

	h.response(w, Ok([]models.Song{song}), http.StatusOK)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	var filters models.Filters
	err := filters.SetQueryData(r)
	if err != nil {
		h.response(w, Error("verse, limit and offset must be int"), http.StatusBadRequest)
	}

	var filter models.Song
	err = filter.SetQueryData(r)
	if err != nil {
		h.response(w, Error("id must be int"), http.StatusBadRequest)
	}

	ctx := logger.NewCtxWithLog(r.Context(), h.log)

	songs, err := h.service.Get(ctx, filter, filters)
	if err != nil {
		h.log.Error(err.Error(), "input",
			slog.Any("filter", filter.AsLogValue()),
			slog.Any("filters", filters.AsLogValue()))

		h.response(w, Error("Can't get songs"), http.StatusInternalServerError)
		return
	}

	h.response(w, Ok(songs), http.StatusOK)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var song models.Song

	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		h.response(w, Error("Can't decode json body"), http.StatusBadRequest)
		return
	}

	ctx := logger.NewCtxWithLog(r.Context(), h.log)

	ok, err := h.service.Update(ctx, song)
	if err != nil {
		h.log.Error(err.Error(), "input", song.AsLogValue())

		h.response(w, Error("Can't update song"), http.StatusInternalServerError)
		return
	}

	if !ok {
		h.response(w, Error("Song not exists"), http.StatusBadRequest)
	}

	h.response(w, Ok(nil), http.StatusOK)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	var song models.Song
	err := song.SetQueryData(r)
	if err != nil {
		h.response(w, Error("id must be int"), http.StatusBadRequest)
	}

	ctx := logger.NewCtxWithLog(r.Context(), h.log)

	ok, err := h.service.Delete(ctx, song.Id)
	if err != nil {
		h.log.Error(err.Error(), "input", song.AsLogValue())

		h.response(w, Error("Can't delete song"), http.StatusInternalServerError)
		return
	}

	if !ok {
		h.response(w, Error("Song not exists"), http.StatusBadRequest)
	}

	h.response(w, Ok(nil), http.StatusNoContent)
}
