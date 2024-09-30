package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/s3nn1k/ef-mob-task/internal/models"
	"github.com/s3nn1k/ef-mob-task/pkg/logger"
)

type Storage struct {
	db PgxPoolIface
}

func (s *Storage) Create(ctx context.Context, song models.Song) (int, error) {
	logger.LogUse(ctx).Debug("Storage.Postgres.Create", "input", song.AsLogValue())

	query := fmt.Sprintf("INSERT INTO %s (song, group_name, text, link, date) VALUES (@song, @group, @text, @link, @date) RETURNING id", table)
	args := pgx.NamedArgs{
		"song":  song.Song,
		"group": song.Group,
		"text":  song.Text,
		"link":  song.Link,
		"date":  song.Date,
	}

	var id int
	err := s.db.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("can't create song in storage: %w", err)
	}

	logger.LogUse(ctx).Debug("Result", slog.Int("id", id))

	return id, nil
}

func (s *Storage) Update(ctx context.Context, song models.Song) (bool, error) {
	logger.LogUse(ctx).Debug("Storage.Postgres.Update", "input", song.AsLogValue())

	query := fmt.Sprintf("UPDATE %s SET song=@song, group_name=@group, text=@text, link=@link, date=@date WHERE id=@id", table)
	args := pgx.NamedArgs{
		"song":  song.Song,
		"group": song.Group,
		"text":  song.Text,
		"link":  song.Link,
		"date":  song.Date,
		"id":    song.Id,
	}

	rows, err := s.db.Exec(ctx, query, args)
	if err != nil {
		return false, fmt.Errorf("can't update song in storage: %w", err)
	}

	res := true
	if rows.RowsAffected() == 0 {
		res = false
	}

	logger.LogUse(ctx).Debug("Result", slog.Bool("updated", res))

	return res, nil

}

func (s *Storage) Delete(ctx context.Context, id int) (bool, error) {
	logger.LogUse(ctx).Debug("Storage.Postgres.Delete", "input", slog.Int("id", id))

	query := fmt.Sprintf("DELETE FROM %s WHERE id=@userId", table)
	args := pgx.NamedArgs{
		"userId": id,
	}

	rows, err := s.db.Exec(ctx, query, args)
	if err != nil {
		return false, fmt.Errorf("can't delete song from storage: %w", err)
	}

	res := true
	if rows.RowsAffected() == 0 {
		res = false
	}

	logger.LogUse(ctx).Debug("Result", slog.Bool("deleted", res))

	return res, nil
}

func (s *Storage) GetAll(ctx context.Context, filters models.GetFilters) ([]models.Song, error) {
	logger.LogUse(ctx).Debug("Storage.Postgres.GetAll", "input", filters.AsLogValue())

	query, args := generateQuery(filters)
	logger.LogUse(ctx).Debug("Generated", slog.Any("query", query), slog.Any("args", args))

	rows, err := s.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	var songs []models.Song
	var logValues []slog.Value
	for rows.Next() {
		var song models.Song

		err := rows.Scan(&song.Id, &song.Song, &song.Group, &song.Text, &song.Link, &song.Date)
		if err != nil {
			return nil, fmt.Errorf("can't get songs from storage: %w", err)
		}

		songs = append(songs, song)
		logValues = append(logValues, song.AsLogValue())
	}

	logger.LogUse(ctx).Debug("Result", slog.Any("songs", logValues))

	return songs, nil
}

// generateQuery generates sql query and []args use given arguments
func generateQuery(filters models.GetFilters) (string, pgx.NamedArgs) {
	query := fmt.Sprintf("SELECT id, song, group_name, text, link, date FROM %s", table)
	var queryArgs []string
	args := pgx.NamedArgs{}

	if filters.Id != 0 {
		queryArgs = append(queryArgs, "id=@id")
		args["id"] = filters.Id
	}

	if filters.Song != "" {
		queryArgs = append(queryArgs, "song=@song")
		args["song"] = filters.Song
	}

	if filters.Group != "" {
		queryArgs = append(queryArgs, "group_name=@group")
		args["group"] = filters.Group
	}

	if filters.Date != "" {
		queryArgs = append(queryArgs, "date=@date")
		args["date"] = filters.Date
	}

	if len(args) > 0 && len(queryArgs) > 0 {
		subQuery := strings.Join(queryArgs, " AND ")

		query += " WHERE " + subQuery
	}

	query += " LIMIT @limit"
	args["limit"] = filters.Limit

	query += " OFFSET @offset"
	args["offset"] = filters.Offset

	return query, args
}
