GO=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go
TAG=1.0.0
BIN=pingdom-exporter
IMAGE=jusbrasil/$(BIN)

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
.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE):$(TAG) .

# Push Docker images to the registry
.PHONY: docker-push
docker-push:
	docker push $(IMAGE):$(TAG)
	docker tag $(IMAGE):$(TAG) $(IMAGE):latest
	docker push $(IMAGE):latest
