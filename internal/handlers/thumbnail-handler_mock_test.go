package handlers

import (
	"context"
	"errors"
	"testing"
	"tumbnail/internal/repository"
	"tumbnail/internal/services"
	"tumbnail/internal/services/mocks"
	"tumbnail/pkg/api"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// MockStream is a mock implementation of the Thumbnail_GetThumbnailServer interface

func TestGetThumbnailSyncSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockThumbnail(ctrl)
	mockStream := new(MockStream)
	repos := &repository.Repository{}
	handler := &Handler{services: services.NewService(repos)}
	handler.services.Thumbnail = mockService

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

func TestGetThumbnailAsyncSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockThumbnail(ctrl)
	mockStream := new(MockStream)
	repos := &repository.Repository{}
	handler := &Handler{services: services.NewService(repos)}
	handler.services.Thumbnail = mockService

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

func TestGetThumbnailSyncError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockThumbnail(ctrl)
	mockStream := new(MockStream)
	repos := &repository.Repository{}
	handler := &Handler{services: services.NewService(repos)}
	handler.services.Thumbnail = mockService

	videoID := "video123"

	mockService.EXPECT().GetThumbnail(videoID).Return(nil, errors.New("not found in Redis")).Times(1)
	mockService.EXPECT().FetchThumbnailFromMicroservice(videoID).Return(nil, errors.New("microservice error")).Times(1)
	mockStream.On("Context").Return(context.Background()).Maybe()

	req := &api.GetThumbnailRequest{VideoID: []string{videoID}, Async: false}
	err := handler.GetThumbnail(req, mockStream)

	assert.Error(t, err)
}

func TestGetThumbnailAsyncError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockThumbnail(ctrl)
	mockStream := new(MockStream)
	repos := &repository.Repository{}
	handler := &Handler{services: services.NewService(repos)}
	handler.services.Thumbnail = mockService

	videoID := "video123"

	mockService.EXPECT().GetThumbnail(videoID).Return(nil, errors.New("not found in Redis")).Times(1)
	mockService.EXPECT().FetchThumbnailFromMicroservice(videoID).Return(nil, errors.New("microservice error")).Times(1)
	mockStream.On("Context").Return(context.Background()).Maybe()

	req := &api.GetThumbnailRequest{VideoID: []string{videoID}, Async: true}
	err := handler.GetThumbnail(req, mockStream)

	assert.Error(t, err)
}
