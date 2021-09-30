.PHONY: build
build: .deps .build

.PHONY: .deps
.deps:
		ls go.mod || go mod init

.PHONY: .build
.build:
		CGO_ENABLED=0 GOOS=linux go build -o bin/tiktok-reporting-api cmd/tiktok-reporting-api/main.go
