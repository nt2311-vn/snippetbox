APP_NAME := web
BIN_DIR  := bin

ifeq ($(OS), Windows_NT)
    EXE := .exe
else
    EXE :=
endif

.PHONY: build run restart

build:
	@echo "Building web app..."
	@go mod tidy
	@go build -o $(BIN_DIR)/$(APP_NAME)$(EXE) ./cmd/web
	@echo "Completed!"

run:
	@echo "Running web app..."
	@$(BIN_DIR)/$(APP_NAME)$(EXE)

restart:
	@echo "Restarting web app..."
	@make build && make run

