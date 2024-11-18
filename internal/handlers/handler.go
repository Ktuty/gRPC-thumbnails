package handlers

import (
	"tumbnail/internal/services"
	"tumbnail/pkg/api"
)

type Handler struct {
	services *services.Service
	api.UnimplementedThumbnailServer
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{services: services}
}
