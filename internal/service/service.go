package service

import (
	"context"
	"strings"

	"github.com/s3nn1k/ef-mob-task/internal/client"
	"github.com/s3nn1k/ef-mob-task/internal/models"
	"github.com/s3nn1k/ef-mob-task/internal/storage"
	"github.com/s3nn1k/ef-mob-task/pkg/logger"
)

// go run github.com/vektra/mockery/v2@v2.45.0 --name=ServiceIface
type ServiceIface interface {
	Create(ctx context.Context, song string, group string) (models.Song, error)
	Update(ctx context.Context, song models.Song) (bool, error)
	GetAll(ctx context.Context, filters models.GetFilters) ([]models.Song, error)
	GetVerses(ctx context.Context, filters models.GetVersesFilters) ([]string, error)
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

func (s *Service) GetVerses(ctx context.Context, filters models.GetVersesFilters) ([]string, error) {
	logger.LogUse(ctx).Debug("Service.GetById", "filters", filters.AsLogValue())

	songs, err := s.storage.GetAll(ctx, models.GetFilters{Limit: 1, Id: filters.Id})
	if err != nil {
		return nil, err
	}

	if len(songs) < 1 {
		return nil, nil
	}

	logger.LogUse(ctx).Debug("Filter song's verses", "input", songs[0].Text)

	verses := filterVerses(songs[0].Text, filters.Limit, filters.Offset)

	logger.LogUse(ctx).Debug("Result", "verses", verses)

	return verses, nil
}

func (s *Service) GetAll(ctx context.Context, filters models.GetFilters) ([]models.Song, error) {
	return s.storage.GetAll(ctx, filters)
}

func (s *Service) Update(ctx context.Context, song models.Song) (bool, error) {
	return s.storage.Update(ctx, song)
}

func (s *Service) Delete(ctx context.Context, id int) (bool, error) {
	return s.storage.Delete(ctx, id)
}

func filterVerses(text string, limit int, offset int) []string {
	verses := strings.Split(text, "\n\n")

	if len(verses) > offset {
		verses = verses[offset:]

		if len(verses) > limit {
			verses = verses[:limit]
		}
	} else {
		return []string{}
	}

	return verses
}
