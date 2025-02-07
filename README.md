# Karma8 Gateway / Storage Backend 

## Как поднять приложение

1. `make compose` - compose сервисов
2. `make migrate` - установка миграций

## Дополнительно
- `make deps` - установка зависимостей
- `make swag` - генерация сваггер документации
  - http://localhost:8080/swagger/gateway/index.html
  - http://localhost:8080/swagger/storage/index.html
- `make migrate` - запуск миграций
- `make migrate-down` - rollback миграций