.PHONY: run build help exe

PROJECT:=project
BATCH:=batch

BIN_DIR:=bin/
BIN_NAME:=main

CMD_DIR:=cmd/
CMD_NAME:=main.go

t:=$(PROJECT)

# help
help:
	@echo "build: build t=*** if t==nil project\nrun: localhost:8000"

# localhost:8000
run:
	@echo "localhost:8000 is started"
	@go run $(CMD_DIR)$(PROJECT)/$(CMD_NAME)

# build
build:
	@go build -o $(BIN_DIR)${t}/$(BIN_NAME) $(CMD_DIR)${t}/$(CMD_NAME)

# execute binary
exe:
	@./$(BIN_DIR)$(PROJECT)/$(BIN_NAME) &

# execute batch binary
batch:
	@./$(BIN_DIR)$(BATCH)/$(BIN_NAME) &
