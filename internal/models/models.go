package models

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
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
type AllFilters struct {
	Limit  int
	Offset int
	Song   string
	Group  string
	Date   string
}

// type SongFilters represents filters that uses for get song with verse
type SongFilters struct {
	Id    int
	Verse int
}

// SetQueryData set's data from request url query to Filters struct
func (a *AllFilters) SetQueryData(r *http.Request) error {
	val := r.URL.Query().Get("limit")
	if val != "" {
		limit, err := strconv.Atoi(val)
		if err != nil {
			return err
		}

		a.Limit = limit
	}

	if a.Limit < 1 {
		a.Limit = 10
	}

	val = r.URL.Query().Get("offset")
	if val != "" {
		offset, err := strconv.Atoi(val)
		if err != nil {
			return err
		}

		a.Offset = offset
	}

	a.Song = r.URL.Query().Get("song")
	a.Group = r.URL.Query().Get("group")
	a.Date = r.URL.Query().Get("date")

	return nil
}

// SetQueryId set's id from request url query to struct
func (s *SongFilters) SetQueryId(r *http.Request) error {
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

// SetQueryVerse set's verse id from request url query to struct
func (s *SongFilters) SetQueryVerse(r *http.Request) error {
	val := r.URL.Query().Get("verse")
	if val != "" {
		verse, err := strconv.Atoi(val)
		if err != nil {
			return err
		}

		s.Verse = verse
	}

	return nil
}

// GetVerse returns verse from text
func (s *Song) GetVerse(id int) string {
	verses := strings.Split(s.Text, "\n\n")

	if len(verses) > id-1 && id >= 1 {
		return verses[id-1]
	}

	return ""
}

// AsLogValue represents AllFilters struct as slog.Value
// Used for logging
func (a *AllFilters) AsLogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("limit", a.Limit),
		slog.Int("offset", a.Offset),
		slog.String("song", a.Song),
		slog.String("group", a.Group),
		slog.String("date", a.Date),
	)
}

// AsLogValue represents SongFilters struct as slog.Value
// Used for logging
func (s *SongFilters) AsLogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("id", s.Id),
		slog.Int("verse", s.Verse),
	)
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
