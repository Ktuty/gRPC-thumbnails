package services

import (
	"fmt"
	"io"
	"net/http"
	"tumbnail/internal/repository"
)

type ThumbnailService struct {
	repo repository.Thumbnail
}

func NewThumbnailService(repo repository.Thumbnail) *ThumbnailService {
	return &ThumbnailService{repo: repo}
}

func (s *ThumbnailService) GetThumbnail(videoID string) ([]byte, error) {
	return s.repo.GetThumbnail(videoID)
}

func (s *ThumbnailService) FetchThumbnailFromMicroservice(videoID string) ([]byte, error) {
	//thumbnailURL := fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", videoID)
	thumbnailURL := fmt.Sprintf("https://img.youtube.com/vi/%s/default.jpg", videoID)
	resp, err := http.Get(thumbnailURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download thumbnail: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func (s *ThumbnailService) SaveThumbnailToRedis(videoID string, thumbnail []byte) error {
	return s.repo.SaveThumbnailToRedis(videoID, thumbnail)
}
