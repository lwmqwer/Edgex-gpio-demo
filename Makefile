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
