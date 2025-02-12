# Karma8 Gateway / Storage Backend 

## Как поднять приложение

1. `make compose` - compose сервисов
2. `make migrate` - установка миграций

# Добавление доп хранилища
1. Создание storage и minio - `make add-storage STORAGE_HOSTNAME=storage7 STORAGE_API_PORT=8087 MINIO_STORAGE_HOST=karma8-minio7 MINIO_PORT=9012 MINIO_CONSOLE_PORT=9013 STORAGE_MINIO_PORT=9000 MINIO_ALIAS=myminio7 MINIO_CONTAINER_NAME=karma8-minio7`
2. [POST] localhost:8080/api/v1/storage - добавление нового хранилища в БД

## Дополнительно
- `make test` - запуск тестов
- `make build` - сборка storage
- `make remove-storage` - отключение доп хранилищ
- `make mock` - генерация моков
- `make swag` - генерация сваггер документации
- `make migrate` - запуск миграций
- `make migrate-down` - rollback миграций


## Docs
- http://localhost:8080/swagger/gateway/index.html
- http://localhost:8080/swagger/storage/index.html
- POSTMAN COLLECTION
- C4