# Image URL to use all building/pushing image targets
# Define registries
STAGING_REGISTRY ?= lwmqwer
IMAGE_NAME ?= edgex-gpio-demo
TAG ?= v0.1.0

IMG ?= ${STAGING_REGISTRY}/${IMAGE_NAME}:${TAG}

all: build

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

build: 
	go build -o bin/device-SAK main.go version.go

tidy:
	go mod tidy

test:
	go test ./... -coverprofile=coverage.out

clean:
	rm -rf bin
	rm coverage.out

docker-build: test ## Build docker image with the manager.
	docker build -t ${IMG} .

docker-push: ## Push docker image with the manager.
	docker push ${IMG}

docker-push-mutiarch:
	docker buildx build --platform linux/arm64,linux/amd64 -t ${IMG} . --push