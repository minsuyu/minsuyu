REPO = github.com/minsuyu/minsuyu

COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null)
BRANCH = $(shell git symbolic-ref --short HEAD)
LATEST_VERSION = $(shell git describe --tags --match "[0-9]*.[0-9]*.[0-9]*" --abbrev=0 2>/dev/null)
VERSION_COMMIT = $(shell git rev-parse --short $(LATEST_VERSION) 2>/dev/null)

# LATEST_VERSION 값이 비어있으면 기본 값을 사용한다.
ifeq ($(LATEST_VERSION), )
LATEST_VERSION = 1.0.0
endif

ifeq ($(COMMIT), $(VERSION_COMMIT))
# 현재 커밋이 최신 버전 커밋인 경우 최신 버전을 사용한다.
SEMANTIC_VERSION = $(LATEST_VERSION)
else
# 현재 커밋이 최신 버전 커밋이 아닌 경우 브랜치 정보를 추가한다.
SEMANTIC_VERSION = $(LATEST_VERSION)-$(BRANCH)
endif


# compile flags
GCFLAGS_DEBUG = all=-N -l

# linker flags
LDFLAGS_VERSION = -X '$(REPO)/internal/version.Version=$(SEMANTIC_VERSION)'
LDFLAGS_COMMIT = -X '$(REPO)/internal/version.Commit=$(COMMIT)'

GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
export CGO_ENABLED=0
export GOTOOLCHAIN=go1.23.12

.PHONY: all
all: \
	insights-scan \
	insights-fsnotify

.PHONY: insights-scan
insights-scan: \
	bin/gs-scan-linux-arm64 \
	bin/gs-scan-linux-amd64

.PHONY: insights-fsnotify
insights-fsnotify: \
	bin/gs-fsnotify-linux-arm64 \
	bin/gs-fsnotify-linux-amd64

bin/gs-scan-linux-arm64:
	@echo "build gs-scan (linux / arm64)"
	@GOOS=linux GOARCH=arm64 \
		go build \
		-o bin/gs-scan-linux-arm64 \
		-ldflags="-s -w $(LDFLAGS_VERSION) $(LDFLAGS_COMMIT)" \
		./cmd/scan

bin/gs-scan-linux-amd64:
	@echo "build gs-scan (linux / amd64)"
	@GOOS=linux GOARCH=amd64 \
		go build \
		-o bin/gs-scan-linux-amd64 \
		-ldflags="-s -w $(LDFLAGS_VERSION) $(LDFLAGS_COMMIT)" \
		./cmd/scan

bin/gs-fsnotify-linux-arm64:
	@echo "build gs-fsnotify (linux / arm64)"
	@GOOS=linux GOARCH=arm64 \
		go build \
		-o bin/gs-fsnotify-linux-arm64 \
		-ldflags="-s -w $(LDFLAGS_VERSION) $(LDFLAGS_COMMIT)" \
		./cmd/fsnotify

bin/gs-fsnotify-linux-amd64:
	@echo "build gs-fsnotify (linux / amd64)"
	@GOOS=linux GOARCH=amd64 \
		go build \
		-o bin/gs-fsnotify-linux-amd64 \
		-ldflags="-s -w $(LDFLAGS_VERSION) $(LDFLAGS_COMMIT)" \
		./cmd/fsnotify

.PHONY: clean
clean:
	@rm -rf bin