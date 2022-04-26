GO=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go
BIN=pingdom-exporter
IMAGE=jusbrasil/$(BIN)
DOCKER_BIN=docker

TAG=$(shell git describe --tags)

.PHONY: build
build:
	$(GO) build -a --ldflags "-X main.VERSION=$(TAG) -w -extldflags '-static'" -tags netgo -o bin/$(BIN) ./cmd/$(BIN)

.PHONY: test
test:
	go vet ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

.PHONY: lint
lint:
	go get -u golang.org/x/lint/golint
	golint ./...

# Build the Docker build stage TARGET
.PHONY: image
image:
	$(DOCKER_BIN) build -t $(IMAGE):$(TAG) .

# Push Docker images to the registry
.PHONY: publish
publish:
	$(DOCKER_BIN) push $(IMAGE):$(TAG)
	$(DOCKER_BIN) tag $(IMAGE):$(TAG) $(IMAGE):latest
	$(DOCKER_BIN) push $(IMAGE):latest
