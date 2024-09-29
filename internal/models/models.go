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

type Filters struct {
	Verse  int
	Limit  int
	Offset int
}

// SetQueryData set's data from request url query to Filters struct
func (f *Filters) SetQueryData(r *http.Request) error {
	val := r.URL.Query().Get("verse")
	if val != "" {
		verse, err := strconv.Atoi(val)
		if err != nil {
			return err
		}

		f.Verse = verse
	}

	val = r.URL.Query().Get("limit")
	if val != "" {
		limit, err := strconv.Atoi(val)
		if err != nil {
			return err
		}

		f.Limit = limit
	}

	if f.Limit < 1 {
		f.Limit = 10
	}

	val = r.URL.Query().Get("offset")
	if val != "" {
		offset, err := strconv.Atoi(val)
		if err != nil {
			return err
		}

		f.Offset = offset
	}

	return nil
}

// SetQueryText set's text data from request url query to Song struct
func (s *Song) SetQueryText(r *http.Request) {
	s.Song = r.URL.Query().Get("song")
	s.Group = r.URL.Query().Get("group")
	s.Text = r.URL.Query().Get("text")
	s.Link = r.URL.Query().Get("link")
	s.Date = r.URL.Query().Get("date")
}

// SetQueryId set's id from request url query to struct
func (s *Song) SetQueryId(r *http.Request) error {
	val := r.URL.Query().Get("id")
	if val != "" {
		id, err := strconv.Atoi(val)
		if err != nil {
			return err
		}

		s.Id = id
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

// AsLogValue represents Filters struct as slog.Value
// Used for logging
func (f *Filters) AsLogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("verse", f.Verse),
		slog.Int("limit", f.Limit),
		slog.Int("offset", f.Offset),
	)
}
