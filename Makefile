# include .env

export

BINARY_DIR=bin/
PROJECT_NAME=bot_server
SERVER=cmd/server/main.go
INFO_DIR=info/

all: compile

compile: | $(BINARY_DIR) $(INFO_DIR)
	@go build -o $(BINARY_DIR)$(PROJECT_NAME) $(SERVER)
	@printf "\033[33mbuild\033[0m\n"

$(BINARY_DIR):
	@mkdir -p bin
	@printf "\033[36mcreate dir binary dir\033[0m\n"

$(INFO_DIR):
	@mkdir -p info
	@printf "\033[36mCreate info dir\033[0m\n"

run:
	@printf "\033[33mRUN\033[0m\n"
	@$(BINARY_DIR)$(PROJECT_NAME)

clean:
	@rm -rf $(INFO_DIR) $(BINARY_DIR)
	@printf "\033[31mdeleted $(INFO_DIR)\033[0m\n\033[31mdeleted $(BINARY_DIR)\033[0m\n\033[33mClean ok\033[0m\n"