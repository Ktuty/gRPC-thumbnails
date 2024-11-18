package handlers

import (
	"log"
	"sync"
	"tumbnail/pkg/api"
)

type DetailedError struct {
	VideoID string
	Error   error
}

func (h *Handler) GetThumbnail(req *api.GetThumbnailRequest, stream api.Thumbnail_GetThumbnailServer) error {
	log.Printf("GetThumbnail called with videoIDs: %v, async: %v", req.VideoID, req.Async)

	if req.Async {
		return h.handleAsync(req, stream)
	} else {
		return h.handleSync(req, stream)
	}
}

func (h *Handler) handleSync(req *api.GetThumbnailRequest, stream api.Thumbnail_GetThumbnailServer) error {
	var thumbnails [][]byte
	var errors []DetailedError

	for _, videoID := range req.VideoID {
		var thumbnail []byte
		var err error
		thumbnail, err = h.services.Thumbnail.GetThumbnail(videoID)
		if err != nil {
			// Если данных нет в Redis, делаем запрос к другому микросервису
			thumbnail, err = h.services.Thumbnail.FetchThumbnailFromMicroservice(videoID)
			if err != nil {
				log.Printf("Error fetching thumbnail from microservice for videoID %s: %v", videoID, err)
				errors = append(errors, DetailedError{VideoID: videoID, Error: err})
				continue
			}
			// Хешируем и сохраняем данные в Redis
			err = h.services.Thumbnail.SaveThumbnailToRedis(videoID, thumbnail)
			if err != nil {
				log.Printf("Error saving thumbnail to Redis for videoID %s: %v", videoID, err)
				errors = append(errors, DetailedError{VideoID: videoID, Error: err})
				continue
			}
		}
		thumbnails = append(thumbnails, thumbnail)
	}

	for _, thumbnail := range thumbnails {
		if err := stream.Send(&api.GetThumbnailResponse{Thumb: thumbnail}); err != nil {
			return err
		}
	}

	if len(errors) > 0 {
		// Возвращаем первую ошибку, если они есть
		return errors[0].Error
	}

	return nil
}

func (h *Handler) handleAsync(req *api.GetThumbnailRequest, stream api.Thumbnail_GetThumbnailServer) error {
	var wg sync.WaitGroup
	results := make(chan []byte, len(req.VideoID))
	errors := make(chan DetailedError, len(req.VideoID))

	for _, videoID := range req.VideoID {
		wg.Add(1)
		go func(videoID string) {
			defer wg.Done()
			var thumbnail []byte
			var err error
			thumbnail, err = h.services.Thumbnail.GetThumbnail(videoID)
			if err != nil {
				// Если данных нет в Redis, делаем запрос к другому микросервису
				thumbnail, err = h.services.Thumbnail.FetchThumbnailFromMicroservice(videoID)
				if err != nil {
					log.Printf("Error fetching thumbnail from microservice for videoID %s: %v", videoID, err)
					errors <- DetailedError{VideoID: videoID, Error: err}
					return
				}
				// Хешируем и сохраняем данные в Redis
				err = h.services.Thumbnail.SaveThumbnailToRedis(videoID, thumbnail)
				if err != nil {
					log.Printf("Error saving thumbnail to Redis for videoID %s: %v", videoID, err)
					errors <- DetailedError{VideoID: videoID, Error: err}
					return
				}
			}
			results <- thumbnail
		}(videoID)
	}

	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	for thumbnail := range results {
		if err := stream.Send(&api.GetThumbnailResponse{Thumb: thumbnail}); err != nil {
			return err
		}
	}

	select {
	case err := <-errors:
		return err.Error
	default:
		return nil
	}
}
