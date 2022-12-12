.PHONY: build clean

CUR_PWD := $(shell pwd)

SERVER_PATH := $(CUR_PWD)
SERVER_BIN := $(CUR_PWD)/web/webserver

export GO111MODULE=on

default: build

build:
	go build -o $(SERVER_BIN) $(SERVER_PATH)/main.go

clean:
	rm $(SERVER_BIN)