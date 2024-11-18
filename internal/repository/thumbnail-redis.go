package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type ThumbnailRedis struct {
	db *redis.Client
}

func NewThumbnailRedis(db *redis.Client) *ThumbnailRedis {
	return &ThumbnailRedis{db}
}

func (t *ThumbnailRedis) GetThumbnail(videoID string) ([]byte, error) {
	return t.db.Get(context.Background(), videoID).Bytes()
}

func (t *ThumbnailRedis) SaveThumbnailToRedis(videoID string, thumbnail []byte) error {
	return t.db.Set(context.Background(), videoID, thumbnail, 0).Err()
}
