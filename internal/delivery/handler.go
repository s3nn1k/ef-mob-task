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
	service service.ServiceIface
}

func NewHandler(l *slog.Logger, s service.ServiceIface) *Handler {
	return &Handler{
		log:     l,
		service: s,
	}
}

// Create creates a new song
// @Summary Create a new song
// @Description Creates a new song with the given details
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body models.Song true "Song details"
// @Success 200 {object} Response "Created song"
// @Failure 400 {object} Response "Invalid input"
// @Failure 500 {object} Response "Failed to create song"
// @Router /songs [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var song models.Song

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		h.response(w, Error("Can't decode json body"), http.StatusBadRequest)
		return
	}

	ctx := logger.NewCtxWithLog(r.Context(), h.log)

	song, err := h.service.Create(ctx, song.Song, song.Group)
	if err != nil {
		h.log.Error(err.Error(), "input", song.AsLogValue())

		h.response(w, Error("Can't create song"), http.StatusInternalServerError)
		return
	}

	h.response(w, Ok([]models.Song{song}), http.StatusOK)
}

// Update updates an existing song
// @Summary Update an existing song
// @Description Updates a song with the given details
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int true "Song Id"
// @Param song body models.Song true "Updated song details"
// @Success 200 {object} Response "Song updated successfully"
// @Failure 400 {object} Response "Invalid input"
// @Failure 404 {object} Response "Song not found"
// @Failure 500 {object} Response "Failed to update song"
// @Router /songs/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var song models.Song

	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		h.response(w, Error("Can't decode json body"), http.StatusBadRequest)
		return
	}

	if err := song.SetQueryId(r); err != nil {
		h.response(w, Error("id must be int"), http.StatusBadRequest)
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
		h.response(w, Error("Song not exists"), http.StatusNotFound)
		return
	}

	h.response(w, Ok(nil), http.StatusOK)
}

// GetAll returns a list of Song's
// @Summary Get all Song's from the storage
// @Description Returns a list of all songs with optional filtering and pagination
// @Tags songs
// @Accept  json
// @Produce  json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param id query int false "Song Id"
// @Param song query string false "Song title"
// @Param group query string false "Group name"
// @Param date query string false "Song release date in format 02.01.2006"
// @Success 200 {object} Response "Array of Song's"
// @Failure 400 {object} Response "Invalid query parameters"
// @Failure 500 {object} Response "Failed to get Song's"
// @Router /songs [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	var filters models.GetFilters

	if err := filters.SetQueryData(r); err != nil {
		h.response(w, Error("limit, offset and id must be int"), http.StatusBadRequest)
		return
	}

	ctx := logger.NewCtxWithLog(r.Context(), h.log)

	songs, err := h.service.GetAll(ctx, filters)
	if err != nil {
		h.log.Error(err.Error(), "input", slog.Any("filters", filters.AsLogValue()))

		h.response(w, Error("Can't get songs"), http.StatusInternalServerError)
		return
	}

	h.response(w, Ok(songs), http.StatusOK)
}

// GetVerses returns paginated verses for a song
// @Summary Get song verses
// @Description Returns paginated verses for the specified song
// @Tags songs
// @Produce  json
// @Param id path int true "Song Id"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} Response "Array of verses"
// @Failure 400 {object} Response "Invalid song Id or pagination parameters"
// @Failure 404 {object} Response "Empty verses response"
// @Failure 500 {object} Response "Failed to get Song's verses"
// @Router /songs/{id} [get]
func (h *Handler) GetVerses(w http.ResponseWriter, r *http.Request) {
	var filters models.GetVersesFilters

	if err := filters.SetQueryId(r); err != nil {
		h.response(w, Error("id must be int"), http.StatusBadRequest)
		return
	}

	if err := filters.SetQueryData(r); err != nil {
		h.response(w, Error("limit and offset must be int"), http.StatusBadRequest)
		return
	}

	ctx := logger.NewCtxWithLog(r.Context(), h.log)

	verses, err := h.service.GetVerses(ctx, filters)
	if err != nil {
		h.log.Error(err.Error(), "input", slog.Any("filters", filters.AsLogValue()))

		h.response(w, Error("Can't get song"), http.StatusInternalServerError)
		return
	}

	if verses == nil {
		h.response(w, Error("Empty verses response"), http.StatusNotFound)
		return
	}

	h.response(w, Ok(verses), http.StatusOK)
}

// Delete deletes a song by Id
// @Summary Delete a song
// @Description Deletes a song with the given Id
// @Tags songs
// @Param id path int true "Song Id"
// @Success 204 {object} Response "Song deleted successfully"
// @Failure 400 {object} Response "Invalid song Id"
// @Failure 404 {object} Response "Song not found"
// @Failure 500 {object} Response "Failed to delete song"
// @Router /songs/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	var filters models.GetVersesFilters

	if err := filters.SetQueryId(r); err != nil {
		h.response(w, Error("id must be int"), http.StatusBadRequest)
		return
	}

	ctx := logger.NewCtxWithLog(r.Context(), h.log)

	ok, err := h.service.Delete(ctx, filters.Id)
	if err != nil {
		h.log.Error(err.Error(), "input", filters.AsLogValue())

		h.response(w, Error("Can't delete song"), http.StatusInternalServerError)
		return
	}

	if !ok {
		h.response(w, Error("Song not exists"), http.StatusNotFound)
		return
	}

	h.response(w, Ok(nil), http.StatusNoContent)
}
