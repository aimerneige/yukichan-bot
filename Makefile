GO						:= go
GO_SOURCES		:= $(shell find . -type f -name "*.go")
TARGET_NAME		:= yukichan-bot
GOOS					?= linux
GOARCH				?= amd64
VERSION				?= v1.0.0
CONF_PATH			?= ./config/config.yaml
DEBUG_LEVEL		?= info

.PHONY: dev run release build fmt clean

dev:
	air --build.cmd "go build -o ./bin/$(TARGET_NAME)-$(GOOS)-$(GOARCH)-$(VERSION) cmd/$(TARGET_NAME)/*" \
		--build.bin "./bin/$(TARGET_NAME)-$(GOOS)-$(GOARCH)-$(VERSION)"

run: build
	./bin/$(TARGET_NAME)-$(GOOS)-$(GOARCH)-$(VERSION)

release: build
	./scripts/release.sh $(VERSION)

build: bin/$(TARGET_NAME)-$(GOOS)-$(GOARCH)-$(VERSION)

bin/$(TARGET_NAME)-$(GOOS)-$(GOARCH)-$(VERSION): $(GO_SOURCES)
	GOOS=$(GOOS) GOARCH=$(GOARCH) \
	$(GO) build -o ./bin/$(TARGET_NAME)-$(GOOS)-$(GOARCH)-$(VERSION) \
	-ldflags "-X main.ConfPath=$(CONF_PATH) -X main.DebugLevel=$(DEBUG_LEVEL)" \
	cmd/$(TARGET_NAME)/*

fmt:
	gofmt -l -w $(GO_SOURCES)

clean:
	-rm -rvf bin/$(TARGET_NAME)-*
