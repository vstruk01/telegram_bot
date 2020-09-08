# include .env

export

BINARY_DIR=bin/
PROJECT_NAME=bot_server
SERVER=cmd/server/main.go
INFO_DIR=info/

all: compile

compile: | $(BINARY_DIR) $(INFO_DIR)
	@go build -o $(BINARY_DIR)$(PROJECT_NAME) $(SERVER)

$(BINARY_DIR):
	@mkdir -p bin

$(INFO_DIR):
	@mkdir -p info

run:
	@$(BINARY_DIR)$(PROJECT_NAME)

clean:
	@rm -rf bin