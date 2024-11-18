package main

import (
	"context"
	"flag"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"strings"
	"tumbnail/internal/client"
	"tumbnail/pkg/api"
)

func main() {
	{
		// Парсинг флага --async
		asyncFlag := flag.Bool("async", false, "Download thumbnails asynchronously")
		flag.Parse()

		// Получение входных данных из терминала
		videoIDsInput := strings.Join(flag.Args(), " ")
		videoIDs := strings.Split(strings.TrimSpace(videoIDsInput), " ")

		// Извлечение videoID из URL
		for i, link := range videoIDs {
			u, err := url.Parse(link)
			if err != nil {
				log.Fatalf("Ошибка при парсинге URL: %v", err)
			}

			// Парсим параметры запроса
			queryParams, err := url.ParseQuery(u.RawQuery)
			if err != nil {
				log.Fatalf("Ошибка при парсинге параметров запроса: %v", err)
			}

			// Извлекаем значение параметра 'v'
			videoID := queryParams.Get("v")
			if videoID == "" {
				log.Fatalf("Параметр 'v' не найден в URL")
			}
			videoIDs[i] = videoID
		}

		if err := initConfig(); err != nil {
			log.Fatalf("error initializing configs: %s", err.Error())
		}

		host := viper.GetString("host")
		port := viper.GetString("port")

		// Создание gRPC клиента
		conn, err := client.NewClient(host + ":" + port)
		if err != nil {
			log.Fatalf("Ошибка при подключении к серверу: %v", err)
		}
		defer conn.Close()

		cl := api.NewThumbnailClient(conn)
		thumbnailClient := client.NewThumbnailClient(cl)

		// Создание запроса
		ctx := context.Background()
		err = thumbnailClient.GetThumbnail(ctx, videoIDs, *asyncFlag)
		if err != nil {
			log.Fatalf("Ошибка при вызове RPC: %v", err)
		}
	}
}
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
