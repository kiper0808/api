# Билд образа с Go
FROM golang:1.22 as builder

WORKDIR /app

## Копируем зависимости и загружаем модули
#COPY go.mod go.sum ./
#RUN go mod download && go mod verify

# Копируем исходный код
COPY . .
RUN go mod download && go mod verify

# Компиляция бинарника
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /app/gateway ./cmd/gateway
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /app/storage ./cmd/storage


# Минимальный образ для запуска
FROM alpine:3.19

RUN apk add --no-cache tzdata

WORKDIR /opt

# Копируем бинарник и шаблоны
COPY --from=builder /app/gateway /opt/gateway
COPY --from=builder /app/storage /opt/storage

# Копируем и настраиваем entrypoint
COPY --from=builder /app/build/entrypoint.sh ./entrypoint.sh
RUN chmod +x ./gateway ./entrypoint.sh
RUN chmod +x ./storage ./entrypoint.sh


# Запуск под некорневым пользователем
RUN adduser -D -u 888 nonroot && chown -R 888:888 .
USER 888

EXPOSE 8080 8081

CMD ["sh", "-c", "./entrypoint.sh"]
