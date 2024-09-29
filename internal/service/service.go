package service

import (
	"context"

	"github.com/s3nn1k/ef-mob-task/internal/client"
	"github.com/s3nn1k/ef-mob-task/internal/models"
	"github.com/s3nn1k/ef-mob-task/internal/storage"
	"github.com/s3nn1k/ef-mob-task/pkg/logger"
)

// go run github.com/vektra/mockery/v2@v2.45.0 --name=ServiceIface
type ServiceIface interface {
	Create(ctx context.Context, song string, group string) (models.Song, error)
	Update(ctx context.Context, song models.Song) (bool, error)
	GetAll(ctx context.Context, filters models.AllFilters) ([]models.Song, error)
	GetById(ctx context.Context, filters models.SongFilters) (models.Song, error)
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

func (s *Service) GetById(ctx context.Context, filters models.SongFilters) (models.Song, error) {
	logger.LogUse(ctx).Debug("Service.GetById", "filters", filters.AsLogValue())

	song, err := s.storage.GetById(ctx, filters.Id)
	if err != nil {
		return models.Song{}, err
	}

	if filters.Verse > 0 {
		logger.LogUse(ctx).Debug("Filter song's text by verse", "input", song.AsLogValue())

		song.Text = song.GetVerse(filters.Verse)

		logger.LogUse(ctx).Debug("Result", "song", song.AsLogValue())
	}

	return song, nil
}

func (s *Service) GetAll(ctx context.Context, filters models.AllFilters) ([]models.Song, error) {
	return s.storage.GetAll(ctx, filters)
}

func (s *Service) Update(ctx context.Context, song models.Song) (bool, error) {
	return s.storage.Update(ctx, song)
}

func (s *Service) Delete(ctx context.Context, id int) (bool, error) {
	return s.storage.Delete(ctx, id)
}
