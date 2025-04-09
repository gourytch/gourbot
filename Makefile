COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date +'%Y-%m-%d %H:%M:%S %Z')
GOURBOT_LDFLAGS := # -ldflags "-X main.ComitHash='$(COMMIT_HASH)' -X main.BuildTime='$(BUILD_TIME)'"

.PHONY: all
all: tidy build

.PHONY: tidy
tidy:
	go mod tidy

# Сборка для локального запуска
.PHONY: build
build:
	go build $(GOURBOT_LDFLAGS) -o bin/gourbot ./cmd/gourbot

# Запуск приложения
.PHONY: run
run: build
	@if [ ! -f .env ]; then \
		echo "Error: .env file is missing. Please create it based on .env-template."; \
		exit 1; \
	fi
	env $(cat .env | xargs) ./bin/gourbot

# Очистка скомпилированных файлов
.PHONY: clean
clean:
	rm -rf bin/

# Запуск тестов
.PHONY: test
test:
	go test ./... -v
