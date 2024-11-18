<div align="center">
  <a href="https://git.io/typing-svg"><img src="https://readme-typing-svg.herokuapp.com?font=Tektur&size=40&duration=4000&color=11FF22&center=true&vCenter=true&width=435&height=100&lines=gRPC-Thumbnails" alt="Typing SVG" /></a>
</div>

# Проект gRPC-thumbnails

Представляет собой систему, которая позволяет скачивать превью (thumbnails) превью изображений с видео YouTube как синхронно, так и асинхронно. 
Этот проект использует gRPC (Google Remote Procedure Call) для взаимодействия между клиентом и сервером, 
что обеспечивает высокую производительность и эффективность при передаче данных.


# Технологии

- Golang
- Redis
- gRPC
- CleanCode

# Подготовка проекта

1. **Клонирование репозитория**:
   ```sh
   git clone https://github.com/Ktuty/gRPC-thumbnails

# Сервер

1. **Установка зависимостей в проекте**:
   ```go
   // из корневой дирректории проетка
   go mod download
   go mod tidy

2. **Пример подключение Redis через Docker**
   ```sh
   docker run --name my-redis -d -p 6379:6379 redis
   
3. **Проверка и изменение файлов конфигурации configs/config.yml:**
   ```yml
   #Example:
   
    port: '8080'
    host: 'localhost'
    db:
      host: 'localhost'
      port: '6379'
      db: '0'

4. **Проверка и изменение файлов конфигурации .env:**
   ```env
   #Example:
   
   DB_PASSWORD=""

5. **Настройка размера картинки в файле .env:**
   ```env
   # Необходимо выбрать размер файла
   
   # Big Image
   IMAGE="https://img.youtube.com/vi/%s/maxresdefault.jpg"

   # Small Image
   IMAGE="https://img.youtube.com/vi/%s/default.jpg"

# Сервер готов к запуску

1. **Запуск сервера:**:
   ```go
   // из корневой дирректории проетка
   go run cmd/server/main.go

# Client 

1. **Запуск клиента:**:
   ```go
   // из корневой дирректории проетка

   // для ассинхронного скачивания 
   go run cmd/server/main.go --async ссыки через пробел

   // для обычного скачивания 
   go run cmd/server/main.go ссыки через пробел

   //Никаких знаков препинания, только пробелы

  ```go
//Example:

// --async
go run cmd/server/main.go --async https://www.youtube.com/watch?v=wEX1_NYoPls https://www.youtube.com/watch?v=xlOjWHWzSYM https://www.youtube.com/watch?v=-gzuWQpQ660 

// usual
go run cmd/server/main.go https://www.youtube.com/watch?v=AxE4TltnvjI https://www.youtube.com/watch?v=a3koSnInh1Y https://www.youtube.com/watch?v=Q9WNn2LRvQk
