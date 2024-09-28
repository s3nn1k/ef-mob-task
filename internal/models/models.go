package models

import (
	"log/slog"
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

// AsLogValue represents Song struct as slog.Value
// Used for logging
func (s *Song) AsLogValue() slog.Value {
	return slog.GroupValue(
		slog.String("song", s.Song),
		slog.String("group", s.Group),
		slog.String("text", s.Text),
		slog.String("link", s.Link),
		slog.String("date", s.Date),
	)
}

// GetVerse returns verse from text
func (s *Song) GetVerse(id int) string {
	verses := strings.Split(s.Text, "/n/n")

	if len(verses) > id-1 {
		return verses[id-1]
	}

	return ""
}
