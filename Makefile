SOURCE_DIR := ./cmd/cli
BUILD_DIR := ./bin
BINARY_NAME := gitecho

all: build

build:
	@GOOS=linux go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SOURCE_DIR)

clean:
	@$(RM) -f $(BUILD_DIR)/$(BINARY_NAME)

.PHONY: all build clean
