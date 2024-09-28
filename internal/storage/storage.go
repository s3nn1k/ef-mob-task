package storage

import (
	"context"

	"github.com/s3nn1k/ef-mob-task/internal/models"
)

type Storage interface {
	Create(ctx context.Context, song models.Song) (int, error)
	Update(ctx context.Context, song models.Song) (bool, error)
	Get(ctx context.Context, filter models.Song, limit int, offset int) ([]models.Song, error)
	Delete(ctx context.Context, id int) (bool, error)
}
