package models

import (
	"log/slog"
	"net/http"
	"strconv"
)

// type Song represents song info
type Song struct {
	Id    int    `json:"id"`
	Song  string `json:"song"`
	Group string `json:"group"`
	Text  string `json:"text"`
	Link  string `json:"link"`
	Date  string `json:"releaseDate"`
}

// type AllFilters represents filters that uses for get library of songs
type GetFilters struct {
	Limit  int
	Offset int
	Id     int
	Song   string
	Group  string
	Date   string
}

// type SongFilters represents filters that uses for get song text with verse
type GetVersesFilters struct {
	Id     int
	Limit  int
	Offset int
}

// SetQueryId set's id from request url query to Song struct
func (s *Song) SetQueryId(r *http.Request) error {
	val := r.PathValue("id")
	if val != "" {
		id, err := strconv.Atoi(val)
		if err != nil {
			return err
		}

		s.Id = id
	}

	return nil
}

// SetQueryData set's data from request url query to AllFilters struct
func (g *GetFilters) SetQueryData(r *http.Request) error {
	val := r.URL.Query().Get("limit")
	if val != "" {
		limit, err := strconv.Atoi(val)
		if err != nil {
			return err
		}

		g.Limit = limit
	}

	if g.Limit < 1 {
		g.Limit = 10
	}

	val = r.URL.Query().Get("offset")
	if val != "" {
		offset, err := strconv.Atoi(val)
		if err != nil {
			return err
		}

		g.Offset = offset
	}

	if g.Offset < 1 {
		g.Offset = 0
	}

	val = r.URL.Query().Get("id")
	if val != "" {
		id, err := strconv.Atoi(val)
		if err != nil {
			return err
		}

		g.Id = id
	}

	g.Song = r.URL.Query().Get("song")
	g.Group = r.URL.Query().Get("group")
	g.Date = r.URL.Query().Get("date")

	return nil
}

// SetQueryId set's id from request url query to GetFilters struct
func (g *GetVersesFilters) SetQueryId(r *http.Request) error {
	val := r.PathValue("id")
	if val != "" {
		id, err := strconv.Atoi(val)
		if err != nil {
			return err
		}

		g.Id = id
	}

	return nil
}

// SetQueryVerse set's data from request url query to GetVersesFilters struct
func (g *GetVersesFilters) SetQueryData(r *http.Request) error {
	val := r.URL.Query().Get("limit")
	if val != "" {
		limit, err := strconv.Atoi(val)
		if err != nil {
			return err
		}

		g.Limit = limit
	}

	if g.Limit < 1 {
		g.Limit = 10
	}

	val = r.URL.Query().Get("offset")
	if val != "" {
		offset, err := strconv.Atoi(val)
		if err != nil {
			return err
		}

		g.Offset = offset
	}

	if g.Offset < 1 {
		g.Offset = 0
	}

	return nil
}

// AsLogValue represents Song struct as slog.Value
// Used for logging
func (s *Song) AsLogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("id", s.Id),
		slog.String("song", s.Song),
		slog.String("group", s.Group),
		slog.String("text", s.Text),
		slog.String("link", s.Link),
		slog.String("date", s.Date),
	)
}

// AsLogValue represents AllFilters struct as slog.Value
// Used for logging
func (g *GetFilters) AsLogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("limit", g.Limit),
		slog.Int("offset", g.Offset),
		slog.Int("id", g.Id),
		slog.String("song", g.Song),
		slog.String("group", g.Group),
		slog.String("date", g.Date),
	)
}

// AsLogValue represents SongFilters struct as slog.Value
// Used for logging
func (g *GetVersesFilters) AsLogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("id", g.Id),
		slog.Int("limit", g.Limit),
		slog.Int("offset", g.Offset),
	)
}
