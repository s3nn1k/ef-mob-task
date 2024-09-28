package service

import (
	"context"

	"github.com/s3nn1k/ef-mob-task/internal/client"
	"github.com/s3nn1k/ef-mob-task/internal/models"
	"github.com/s3nn1k/ef-mob-task/internal/storage"
)

// go run github.com/vektra/mockery/v2@v2.45.0 --name=Service
type ServiceIface interface {
	Create(ctx context.Context, song string, group string) (models.Song, error)
	Update(ctx context.Context, song models.Song) (bool, error)
	Get(ctx context.Context, filter models.Song, limit int, offset int) ([]models.Song, error)
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

func (s *Service) Get(ctx context.Context, filter models.Song, limit int, offset int) ([]models.Song, error) {
	return s.storage.Get(ctx, filter, limit, offset)
}

func (s *Service) Delete(ctx context.Context, id int) (bool, error) {
	return s.storage.Delete(ctx, id)
}
