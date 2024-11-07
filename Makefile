# 定义变量
APP_DIR := my-app
PROJECT_NAME=TalentRank
MAIN_FILE=main.go
PKG := "github.com/acd19ml/$(PROJECT_NAME)"
BUILD_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_COMMIT := ${shell git rev-parse HEAD}
BUILD_TIME := ${shell date '+%Y-%m-%d %H:%M:%S'}
BUILD_GO_VERSION := $(shell go version | grep -o  'go[0-9].[0-9].*')
VERSION_PATH := "${PKG}/version"

# 定义命令
init:
	cd $(APP_DIR) && npm init -y

install:
	cd $(APP_DIR) && npm install && npm install antd --save

startf:
	cd $(APP_DIR) && npm start

runf: init install startf

run:
	go run main.go start

dep: ## Get the dependencies
	@go mod tidy

build: dep ## Build the binary file
	@go build -ldflags "-s -w" -ldflags "-X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" -o dist/TalentRank $(MAIN_FILE)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
