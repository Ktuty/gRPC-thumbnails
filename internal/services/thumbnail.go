package services

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"tumbnail/internal/repository"

	"github.com/joho/godotenv"
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

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	if os.Getenv("IMAGE") == "" {
		log.Fatal("IMAGE URL is not set in environment variables")
	}

	thumbnailURL := fmt.Sprintf(os.Getenv("IMAGE"), videoID)
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
