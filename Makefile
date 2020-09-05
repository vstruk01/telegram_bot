# include .env

export

BINARY_DIR=bin/
PROJECT_NAME=t_e_bot
MAIN=main

all: compile

compile: | $(BINARY_DIR)
	@go build -o $(BINARY_DIR)$(PROJECT_NAME) main.go

$(BINARY_DIR):
	@mkdir -p bin

run:
	@$(BINARY_DIR)$(PROJECT_NAME)

clean:
	@rm -rf 