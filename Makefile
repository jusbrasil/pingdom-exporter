GO=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go
TAG=1.0.0
BIN=pingdom_exporter
IMAGE=jusbrasil/$(BIN)

.PHONY: build
build:
	$(GO) build -a --ldflags "-X main.VERSION=$(TAG) -w -extldflags '-static'" -tags netgo -o bin/$(BIN) ./cmd/$(BIN)

.PHONY: image
image: build
	docker build -t $(IMAGE):$(TAG) .

.PHONY: push
push: image
	docker push $(IMAGE):$(TAG)

.PHONY: push-latest
push-latest: image
	docker tag $(IMAGE):$(TAG) $(IMAGE):latest
	docker push $(IMAGE):latest

.PHONY: clean
clean:
	rm -Rf bin/ cover*

.PHONY: test
test:
	ls

.PHONY: cover
cover:
	ls
