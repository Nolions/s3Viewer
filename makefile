APP_NAME := s3Viewer
CMD_PATH := ./cmd/app
OUTPUT_DIR := build

.PHONY: windows-amd64 linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 clean

default: build

run:
	@echo "run..."
	go run $(CMD_PATH)/main.go

compile:
	@echo "Building by OS..."
	go build -o $(OUTPUT_DIR)/$(APP_NAME).exe $(CMD_PATH)

windows-amd64:
	@echo "Building for Windows AMD64..."
	go env -w GOOS=windows
	go env -w GOARCH=amd64
	go build -o $(OUTPUT_DIR)/$(APP_NAME)-windows-amd64.exe $(CMD_PATH)
	go env -u GOOS
	go env -u GOARCH

linux-amd64:
	@echo "Building for Linux AMD64..."
	go env -w GOOS=linux
	go env -w GOARCH=amd64
	go build -o $(OUTPUT_DIR)/$(APP_NAME)-linux-amd64 $(CMD_PATH)
	go env -u GOOS
	go env -u GOARCH

linux-arm64:
	@echo "Building for Linux ARM64..."
	go env -w GOOS=linux
	go env -w GOARCH=arm64
	go build -o $(OUTPUT_DIR)/$(APP_NAME)-linux-arm64 $(CMD_PATH)
	go env -u GOOS
	go env -u GOARCH

darwin-amd64:
	@echo "Building for macOS AMD64..."
	go env -w GOOS=darwin
	go env -w GOARCH=amd64
	go build -o $(OUTPUT_DIR)/$(APP_NAME)-darwin-amd64 $(CMD_PATH)
	go env -u GOOS
	go env -u GOARCH

darwin-arm64:
	@echo "Building for macOS ARM64..."
	go env -w GOOS=darwin
	go env -w GOARCH=arm64
	go build -o $(OUTPUT_DIR)/$(APP_NAME)-darwin-arm64 $(CMD_PATH)
	go env -u GOOS
	go env -u GOARCH

clean:
	rm -rf $(OUTPUT_DIR)

runMinIO:
	docker run -d --name minio -p 9000:9000 \
	-p 9001:9001 \
	-e MINIO_ROOT_USER=admin \
	-e MINIO_ROOT_PASSWORD=admin12345 \
	-v ~/volume/minio/data:/data \
	-v ~/volume/minio/config:/root/.minio \
	quay.io/minio/minio server /data --console-address ":9001"
