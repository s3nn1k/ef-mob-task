package storage

import (
	"context"

	"github.com/s3nn1k/ef-mob-task/internal/models"
)

// go run github.com/vektra/mockery/v2@v2.45.0 --name=Storage
type Storage interface {
	Create(ctx context.Context, song models.Song) (int, error)
	Update(ctx context.Context, song models.Song) (bool, error)
	GetAll(ctx context.Context, filters models.GetFilters) ([]models.Song, error)
	Delete(ctx context.Context, id int) (bool, error)
}
