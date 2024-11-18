package services

import (
	"tumbnail/internal/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type Thumbnail interface {
	GetThumbnail(videoID string) ([]byte, error)
	FetchThumbnailFromMicroservice(videoID string) ([]byte, error)
	SaveThumbnailToRedis(videoID string, thumbnail []byte) error
}

type Service struct {
	Thumbnail
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Thumbnail: NewThumbnailService(repos.Thumbnail),
	}
}
