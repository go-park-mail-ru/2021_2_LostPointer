LDFLAGS   = \
			-X $(GOPREFIX).Version=$(VERSION) \
			-X $(GOPREFIX).Branch=$(BRANCH) \
			-X $(GOPREFIX).Revision=$(REVISION)

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

.PHONY: build
build:
ifndef TARGET
	@echo 'build target is not defined'
else
	go build $(GOTAGS) \
		-ldflags '$(LDFLAGS)' \
		-o bin/$(TARGET) \
		./cmd/$(TARGET)
endif

.PHONY: run
run:
	docker-compose up -d --build
