package repository

import "github.com/go-redis/redis/v8"

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type Thumbnail interface {
	GetThumbnail(videoID string) ([]byte, error)
	SaveThumbnailToRedis(videoID string, thumbnail []byte) error
}

type Repository struct {
	Thumbnail
}

func NewRepository(db *redis.Client) *Repository {
	return &Repository{
		Thumbnail: NewThumbnailRedis(db),
	}
}
