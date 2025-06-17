# Define vars
BINARY_NAME := ./igor
SRC_DIR := ./src
MAIN := $(SRC_DIR)/main.go
VENV := ./.venv
# Run the project
run:
	OLLAMA_HOST=http://mother.local:11434 go run ./src

venv:
	@if [ ! -d .venv ]; then \
		echo "Creating virtual environment..."; \
		python3.11 -m venv $(VENV); \
	fi
	@$(VENV)/bin/pip install piper-tts --no-deps piper-phonemize-cross onnxruntime numpy

# Build the binary
build:
	go mod tidy
	go build -o $(BINARY_NAME) $(SRC_DIR)

# Clean build artifacts
clean:
	rm -rf $(VENV)
	rm -f $(BINARY_NAME)
	rm -rf go.sum
