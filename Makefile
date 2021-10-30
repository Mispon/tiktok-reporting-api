.PHONY: build
build: .deps .build .copy

.PHONY: .deps
.deps:
		ls go.mod || go mod init

.PHONY: .build
.build:
		CGO_ENABLED=0 GOOS=linux go build -o bin/tiktok-reporting-api cmd/tiktok-reporting-api/main.go

.PHONY: .copy
.copy:
	cp configs/credentials.json bin/
	cp configs/advert_ids.txt bin/
