package handlers

import (
	"context"
	"errors"
	"log"
	"sync"
	"testing"
	"tumbnail/internal/repository"
	"tumbnail/internal/services"
	"tumbnail/internal/services/mocks"
	"tumbnail/pkg/api"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

func (h *Handler) GetThumbnailTest(req *api.GetThumbnailRequest, stream api.Thumbnail_GetThumbnailServer) error {
	log.Printf("GetThumbnail called with videoIDs: %v, async: %v", req.VideoID, req.Async)

	if req.Async {
		return h.handleAsync(req, stream)
	} else {
		return h.handleSync(req, stream)
	}
}

func (h *Handler) handleSyncTest(req *api.GetThumbnailRequest, stream api.Thumbnail_GetThumbnailServer) error {
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

func (h *Handler) handleAsyncTest(req *api.GetThumbnailRequest, stream api.Thumbnail_GetThumbnailServer) error {
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

// MockStream is a mock implementation of the Thumbnail_GetThumbnailServer interface
type MockStream struct {
	mock.Mock
	grpc.ServerStream
}

func (m *MockStream) Send(response *api.GetThumbnailResponse) error {
	args := m.Called(response)
	return args.Error(0)
}

func (m *MockStream) Context() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

func TestHandleSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockThumbnail(ctrl)
	mockStream := new(MockStream)
	repos := &repository.Repository{}                         // Создаем пустой репозиторий
	handler := &Handler{services: services.NewService(repos)} // Используем services.NewService
	handler.services.Thumbnail = mockService                  // Заменяем реальный сервис моком

	videoID := "video123"
	thumbnail := []byte("thumbnail_data")

	mockService.EXPECT().GetThumbnail(videoID).Return(nil, errors.New("not found in Redis")).Times(1)
	mockService.EXPECT().FetchThumbnailFromMicroservice(videoID).Return(thumbnail, nil).Times(1)
	mockService.EXPECT().SaveThumbnailToRedis(videoID, thumbnail).Return(nil).Times(1)
	mockStream.On("Send", &api.GetThumbnailResponse{Thumb: thumbnail}).Return(nil).Once()
	mockStream.On("Context").Return(context.Background()).Maybe()

	req := &api.GetThumbnailRequest{VideoID: []string{videoID}, Async: false}
	err := handler.GetThumbnail(req, mockStream)

	assert.NoError(t, err)
	mockStream.AssertExpectations(t)
}

func TestHandleAsync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockThumbnail(ctrl)
	mockStream := new(MockStream)
	repos := &repository.Repository{}                         // Создаем пустой репозиторий
	handler := &Handler{services: services.NewService(repos)} // Используем services.NewService
	handler.services.Thumbnail = mockService                  // Заменяем реальный сервис моком

	videoID := "video123"
	thumbnail := []byte("thumbnail_data")

	mockService.EXPECT().GetThumbnail(videoID).Return(nil, errors.New("not found in Redis")).Times(1)
	mockService.EXPECT().FetchThumbnailFromMicroservice(videoID).Return(thumbnail, nil).Times(1)
	mockService.EXPECT().SaveThumbnailToRedis(videoID, thumbnail).Return(nil).Times(1)
	mockStream.On("Send", &api.GetThumbnailResponse{Thumb: thumbnail}).Return(nil).Once()
	mockStream.On("Context").Return(context.Background()).Maybe()

	req := &api.GetThumbnailRequest{VideoID: []string{videoID}, Async: true}
	err := handler.GetThumbnail(req, mockStream)

	assert.NoError(t, err)
	mockStream.AssertExpectations(t)
}

func TestHandleSyncError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockThumbnail(ctrl)
	mockStream := new(MockStream)
	repos := &repository.Repository{}                         // Создаем пустой репозиторий
	handler := &Handler{services: services.NewService(repos)} // Используем services.NewService
	handler.services.Thumbnail = mockService                  // Заменяем реальный сервис моком

	videoID := "video123"

	mockService.EXPECT().GetThumbnail(videoID).Return(nil, errors.New("not found in Redis")).Times(1)
	mockService.EXPECT().FetchThumbnailFromMicroservice(videoID).Return(nil, errors.New("microservice error")).Times(1)
	mockStream.On("Context").Return(context.Background()).Maybe()

	req := &api.GetThumbnailRequest{VideoID: []string{videoID}, Async: false}
	err := handler.GetThumbnail(req, mockStream)

	assert.Error(t, err)
}

func TestHandleAsyncError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockThumbnail(ctrl)
	mockStream := new(MockStream)
	repos := &repository.Repository{}                         // Создаем пустой репозиторий
	handler := &Handler{services: services.NewService(repos)} // Используем services.NewService
	handler.services.Thumbnail = mockService                  // Заменяем реальный сервис моком

	videoID := "video123"

	mockService.EXPECT().GetThumbnail(videoID).Return(nil, errors.New("not found in Redis")).Times(1)
	mockService.EXPECT().FetchThumbnailFromMicroservice(videoID).Return(nil, errors.New("microservice error")).Times(1)
	mockStream.On("Context").Return(context.Background()).Maybe()

	req := &api.GetThumbnailRequest{VideoID: []string{videoID}, Async: true}
	err := handler.GetThumbnail(req, mockStream)

	assert.Error(t, err)

}
