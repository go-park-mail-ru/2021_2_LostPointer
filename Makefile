.PHONY: lint
ifeq ($(BRANCH),master)
lint:
	@echo "No need run linters on branch $(BRANCH)"
else
lint:
	golangci-lint run -c golangci.yml ./...
endif

.PHONY: test
test:
	go test ./... -covermode=atomic -v -race

.PHONY: all
all: lint test

.PHONY: run
run:
	docker-compose up -d --build
