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

.PHONY: tests
tests:
	go test -coverpkg=./... -coverprofile cover.out.tmp ./...
	cat cover.out.tmp grep -v "monitoring" | grep -v "easyjson" | grep -v "mock_*" | grep -v ".pb.go" | grep -v ".pb" | grep -v "middleware.go" | grep -v "/cmd*"> cover.out
	go tool cover -func cover.out

.PHONY: build_local
build_local:
	docker-compose -f docker-compose.yml up --build -d

.PHONY: build_prod
build_prod:
	docker-compose -f docker-compose-cd.yml up --build -d
