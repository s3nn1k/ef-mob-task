package service

import (
	"context"
	"log/slog"

	"github.com/s3nn1k/ef-mob-task/internal/client"
	"github.com/s3nn1k/ef-mob-task/internal/models"
	"github.com/s3nn1k/ef-mob-task/internal/storage"
	"github.com/s3nn1k/ef-mob-task/pkg/logger"
)

// go run github.com/vektra/mockery/v2@v2.45.0 --name=Service
type ServiceIface interface {
	Create(ctx context.Context, song string, group string) (models.Song, error)
	Update(ctx context.Context, song models.Song) (bool, error)
	Get(ctx context.Context, filter models.Song, filters models.Filters) ([]models.Song, error)
	Delete(ctx context.Context, id int) (bool, error)
}

type Service struct {
	storage storage.Storage
	client  client.ClientIface
}

func New(s storage.Storage, c client.ClientIface) ServiceIface {
	return &Service{
		storage: s,
		client:  c,
	}
}

func (s *Service) Create(ctx context.Context, song string, group string) (models.Song, error) {
	res, err := s.client.GetDetail(ctx, song, group)
	if err != nil {
		return models.Song{}, err
	}

	id, err := s.storage.Create(ctx, res)
	if err != nil {
		return models.Song{}, err
	}

	res.Id = id

	return res, nil
}

func (s *Service) Update(ctx context.Context, song models.Song) (bool, error) {
	return s.storage.Update(ctx, song)
}

func (s *Service) Get(ctx context.Context, filter models.Song, filters models.Filters) ([]models.Song, error) {
	logger.LogUse(ctx).Debug("Service.Get", "filters", filters.AsLogValue())

	songs, err := s.storage.Get(ctx, filter, filters.Limit, filters.Offset)
	if err != nil {
		return nil, err
	}

	if filters.Verse > 0 {
		logger.LogUse(ctx).Debug("Filter song's text by verse", "input", songs)

		for index, song := range songs {
			song.Text = song.GetVerse(filters.Verse)

			songs[index] = song
		}

		logger.LogUse(ctx).Debug("Result", slog.Any("songs", songs))
	}

	return songs, nil
}

func (s *Service) Delete(ctx context.Context, id int) (bool, error) {
	return s.storage.Delete(ctx, id)
}
