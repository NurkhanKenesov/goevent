FROM golang:1.25.3-alpine
WORKDIR /app

# Копируем зависимости и скачиваем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Команда по умолчанию запускает сервер
CMD ["go", "run", "cmd/app/main.go"]

