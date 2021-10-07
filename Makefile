.PHONY: run build help print

PROJECT:=project

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

print:
	@echo ${t}
