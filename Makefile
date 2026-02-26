# Установка зависимостей
deps: .install-deps

# Генерация
gen: .install-deps .mockgen

# Линт проекта
lint-all: .lint-changes .lint-full

# Запуск docker
up: .docker-up

run:
	godotenv -f dev/.env go run cmd/main.go


# Установка зависимостей
.install-deps:
	go install tool


.mockgen:
	@go generate -run minimock ./...



# Проверяем не сделали ли мы обратно-несовместимые изменения контракта по сравнению с develop веткой для разработки
.lint-changes:
	golangci-lint run \
		--new-from-rev=origin/main \
		--config=.golangci.yml \
		./...


# Линт всех файлов
.lint-full:
	golangci-lint run \
		--config=.golangci.yml \
		./...


# Докер
.docker-up:
	@cd dev && docker compose up
