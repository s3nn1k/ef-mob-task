package models

import "log/slog"

// type Song represents song info
type Song struct {
	Id    int
	Song  string
	Group string
	Text  string
	Link  string
	Date  string
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
