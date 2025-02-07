-include .env

export GO111MODULE=on
export CGO_ENABLED=0

# colors
GREEN=\033[1;32m
PURPLE=\033[1;35m
NC=\033[0m

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

# lint
LINTER_VERSION=1.57.2
lint:
	@echo 'run golangci lint'
	@if ! $(GOPATH)/bin/golangci-lint --version | grep -q $(LINTER_VERSION); \
		then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v$(LINTER_VERSION); fi;
	$(GOPATH)/bin/golangci-lint run --out-format=tab -v --whole-files
	@echo
