BINARY := tokcount
BUILD_DIR := bin

.PHONY: build run test tidy clean

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY) .

run:
	go run . .

test:
	go test . ./cmd/tokcount ./internal/cli ./internal/count ./internal/ignore ./internal/output ./internal/tokenizer

tidy:
	go mod tidy

clean:
	trash $(BUILD_DIR) 2>/dev/null || true
