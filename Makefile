-include .env

export GO111MODULE=on
export CGO_ENABLED=0

# colors
GREEN=\033[1;32m
PURPLE=\033[1;35m
NC=\033[0m

# make mocks
mock:
	@echo 'generate mocks'
	go generate ./internal/gateway/service
	go generate ./internal/gateway/service/file_storage

# make tests
test:
	@echo 'run tests'
	go test -v -cover ./internal/gateway/... ./internal/storage/...

# compose deps
compose:
	@echo 'compose deps'
	docker-compose -f docker-compose.yaml up -d

# down deps
compose-down:
	@echo 'compose deps'
	docker-compose -f docker-compose.yaml down

# generate swagger
swag:
	@echo 'generation swagger docs'
	swag init --parseDependency -g handler.go -dir internal/gateway/api/http/v1 --instanceName gateway
	swag init --parseDependency -g handler.go -dir internal/storage/api/http/v1 --instanceName storage

# migrate
migrate:
	@echo "\n${GREEN}UP MIGRATE DB${NC}\n"
	@docker run -e INSTALL_MYSQL=true --rm -it \
		  -v ./dev/liquibase/changelogs/karma8/changelog.sql:/liquibase/changelog/changelog.sql \
		  --env-file dev/liquibase/liquibase.docker.karma8.env \
		  liquibase/liquibase update --log-level info

# migrate-rollback
migrate-down:
	@echo "\n${PURPLE}ROLLBACK MIGRATE DB${NC}\n"
	@docker run -e INSTALL_MYSQL=true --rm -it \
		  -v ./dev/liquibase/changelogs/karma8/changelog.sql:/liquibase/changelog/changelog.sql \
		  --env-file dev/liquibase/liquibase.docker.karma8.env \
		  liquibase/liquibase rollback-count --count=1

# build storage
build-storage:
	@echo "\n${GREEN}BUILD STORAGE${NC}\n"
	docker build -f ./build/Dockerfile.storage -t karma8-storage . --no-cache

# add storage
.PHONY: add-storage
add-storage:
	# Проверяем, что все необходимые переменные заданы
	@if [ -z "$(STORAGE_HOSTNAME)" ]; then \
		echo "Ошибка: STORAGE_HOSTNAME не задано!"; \
		exit 1; \
	fi
	@if [ -z "$(STORAGE_API_PORT)" ]; then \
		echo "Ошибка: STORAGE_API_PORT не задано!"; \
		exit 1; \
	fi
	@if [ -z "$(MINIO_STORAGE_HOST)" ]; then \
		echo "Ошибка: MINIO_STORAGE_HOST не задано!"; \
		exit 1; \
	fi
	@if [ -z "$(MINIO_PORT)" ]; then \
		echo "Ошибка: MINIO_PORT не задано!"; \
		exit 1; \
	fi
	@if [ -z "$(MINIO_CONSOLE_PORT)" ]; then \
		echo "Ошибка: MINIO_CONSOLE_PORT не задано!"; \
		exit 1; \
	fi
	@if [ -z "$(MINIO_ALIAS)" ]; then \
    	echo "Ошибка: MINIO_ALIAS не задано!"; \
    	exit 1; \
    fi

	# Создаем volume для MinIO
	docker volume create $(MINIO_STORAGE_HOST)_data

	# Запускаем docker-compose с переданными переменными
	MINIO_CONTAINER_NAME="$(MINIO_STORAGE_HOST)" \
	STORAGE_CONTAINER_NAME="karma8-$(STORAGE_HOSTNAME)" \
	MINIO_STORAGE_HOST=$(MINIO_STORAGE_HOST) \
	MINIO_PORT=$(MINIO_PORT) \
	MINIO_ALIAS=$(MINIO_ALIAS) \
	MINIO_CONSOLE_PORT=$(MINIO_CONSOLE_PORT) \
	STORAGE_HOSTNAME=$(STORAGE_HOSTNAME) \
	STORAGE_API_PORT=$(STORAGE_API_PORT) \
	STORAGE_MINIO_PORT=$(STORAGE_MINIO_PORT) \
	docker-compose -f add-storage-dc.yaml up -d

# remove storage
.PHONY: remove-storage
remove-storage:
	docker-compose -f add-storage-dc.yaml down