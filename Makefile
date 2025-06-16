# Define vars
BINARY_NAME := ./igor
SRC_DIR := ./src
MAIN := $(SRC_DIR)/main.go

# Run the project
run:
	go run ./src

# Build the binary
build:
	go mod tidy
	go build -o $(BINARY_NAME) $(SRC_DIR)

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf go.sum
