package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"tumbnail/pkg/api"
)

type ThumbnailClient struct {
	client api.ThumbnailClient
}

func NewThumbnailClient(client api.ThumbnailClient) *ThumbnailClient {
	return &ThumbnailClient{client: client}
}

func (c *ThumbnailClient) GetThumbnail(ctx context.Context, videoIDs []string, asyncFlag bool) error {
	log.Printf("\nGetThumbnail called with --async: %v\nvideoIDs: %v\n\n", asyncFlag, videoIDs)
	stream, err := c.client.GetThumbnail(ctx, &api.GetThumbnailRequest{VideoID: videoIDs, Async: asyncFlag})
	if err != nil {
		return err
	}

	// Создание директории для вывода, если она не существует
	outputDir := "cmd/outputs"
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("ошибка при создании директории: %v", err)
	}

	// Индекс для именования файлов
	index := 0

	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("ошибка при получении ответа: %v", err)
		}

		// Запись миниатюры в файл
		filename := filepath.Join(outputDir, fmt.Sprintf("%s.jpg", videoIDs[index]))
		err = os.WriteFile(filename, resp.Thumb, 0644)
		if err != nil {
			return fmt.Errorf("ошибка при записи файла: %v", err)
		}
		fmt.Printf("Миниатюра сохранена в файл: %s\n", filename)
		index++
	}

	return nil
}
