.DEFAULT_GOAL := help

files := $(shell find . -path ./vendor -prune -path ./pb -prune -o -name '*.go' -print)
pkgs := $(shell go list ./... | grep -v /vendor/ )

git_rev := $(shell git rev-parse --short HEAD)
git_tag := $(shell git tag --points-at=$(git_rev))
release_date := $(shell date +%d-%m-%Y)
latest_git_tag := $(shell git for-each-ref --format="%(tag)" --sort=-taggerdate refs/tags | head -1)
ifneq ($(strip $(latest_git_tag)),)
	latest_git_rev := $(shell git rev-list --abbrev-commit -n 1 $(latest_git_tag))
endif
version := $(if $(git_tag),$(git_tag),dev@$(git_rev))
build_time := $(shell date -u)
ldflags := -X "main.Version=$(version)" -X "main.BuildTime=$(build_time)"

cwd := $(shell pwd)
build_dir := $(cwd)/build/bin
dist_dir := $(cwd)/dist

# Define cross compiling targets
os := $(shell uname)
ifeq ("$(os)", "Linux")
	target_os = linux
	cross_os = darwin
else ifeq ("$(os)", "Darwin")
	target_os = darwin
	cross_os = linux
endif

# Debug purpose
print-%: ; @echo $*=$($*)

# Define cross compiling targets
os := $(shell uname)

.PHONY: check-os
check-os:
ifndef target_os
	$(error Unsupported platform: ${os})
endif

.PHONY: format
format:
	@echo "===== format"
	@goimports -w $(files)
	@sync

unformatted = $(shell goimports -l $(files))

.PHONY: check-format
check-format:
	@echo "===== check formatting"
ifneq "$(unformatted)" ""
	@echo "needs formatting:"
	@echo "$(unformatted)" | tr ' ' '\n'
	$(error run 'make format')
endif

.PHONY: vet
vet:
	@echo "===== vet"
	@go vet $(pkgs)

.PHONY: lint
lint:
	@echo "===== lint"
	@for pkg in $(pkgs); do \
		golint -set_exit_status $$pkg || exit 1 ; \
	done;

.PHONY: check
check: check-os check-format vet lint

.PHONY: setup
setup:
	@echo "===== setup"
	go get -v golang.org/x/lint/golint
	go get golang.org/x/tools/cmd/goimports
	@which golint > /dev/null || (echo 'ERROR: unable to find golint'; exit 1)
	@which goimports > /dev/null || (echo 'ERROR: unable to find goimports'; exit 1)

.PHONY: build
build: setup check ## Build artifact in 'build' directory for the current platform
	@echo "===== build"
	GOOS=${target_os} GOARCH=amd64 go build -ldflags '-s $(ldflags)' -o ${build_dir}/go-hello -v

.PHONY: clean
clean:
	@echo "===== clean"
	rm -rf build

.PHONY: test
test: build ## Run all tests from source files
	@echo "===== run tests"
	go test -race $(pkgs)

image := nocquidant/go-hello
DOCKER_USER  ?= tobeset-here
DOCKER_PASS ?= tobeset-cicd

.PHONY: docker-check
docker-check: 
	@echo "===== docker check env."
	@which docker > /dev/null

.PHONY: docker-image
docker-image: docker-check ## Build Docker image
	@echo "===== build docker image"
	docker build -t $(image):latest .

.PHONY: docker-release
docker-release: docker-image  ## Push Docker image
	@echo "===== tag and push docker image"
ifeq ($(strip $(git_tag)),)
	@echo "no tag on $(git_rev), skipping docker release"
else
	@echo "releasing $(image):$(git_tag)"
	@docker login -u $(DOCKER_USER) -p $(DOCKER_PASS)
	docker tag $(image):latest $(image):$(git_tag)
	docker push $(image):$(git_tag)
	@if [ "$(git_rev)" = "$(latest_git_rev)" ]; then \
		echo "updating latest image"; \
		echo docker push $(image):latest ; \
	fi;
endif

.PHONY: help
help: ## Displays this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)