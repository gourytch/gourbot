
COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date +'%Y-%m-%dT%H:%M:%S')
VERSION=$(shell git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')
VPKG := gourbot/internal
GOURBOT_LDFLAGS := -ldflags "-X ${VPKG}/version.Version='${VERSION}' -X ${VPKG}/version.CommitHash='${COMMIT_HASH}' -X ${VPKG}/version.BuildTime='${BUILD_TIME}'"

.PHONY: all
all: tidy build

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: build
build: ./bin/gourbot

./bin/gourbot::
	go build  $(GOURBOT_LDFLAGS) -o ./bin/gourbot ./cmd/gourbot

.PHONY: run
run: build
	@if [ ! -f .env ]; then \
		echo "Error: .env file is missing. Please create it based on .env-template."; \
		exit 1; \
	fi
#	env $(cat .env | xargs)
	./bin/gourbot

.PHONY: clean
clean:
	rm -rf bin/ log/ || true

# Запуск тестов
.PHONY: test
test:
	go test ./... -v
