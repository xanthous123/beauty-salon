.PHONY: test cover build clean help

# Переменные
BINARY_NAME=beauty-salon
COVERAGE_FILE=coverage.out

## test: запуск всех юнитов-тестов проекта
test:
	go test -v ./...

## cover: запуск тестов с генерацией отчета о покрытии и открытием его в браузере
cover:
	go test -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -func=$(COVERAGE_FILE)
	@echo "Opening HTML report..."
	go tool cover -html=$(COVERAGE_FILE)

## race: запуск тестов с проверкой на состояние гонки (race detector)
race:
	go test -race ./...

## build: сборка бинарного файла приложения
build:
	go build -o $(BINARY_NAME) ./cmd/main.go

## clean: очистка временных файлов и бинарника
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(COVERAGE_FILE)

## help: справка по командам
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: build-secure

build-secure:
	go install mvdan.cc/garble@latest
	# Сборка обфусцированного файла
	garble -literals -tiny -seed=random build -o bin/beauty-salon-secure ./cmd/main.go