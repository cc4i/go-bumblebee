
# Image URL to use all building/pushing
IMG ?= gate:latest
# Docker Hub username / password
DOCKER_PASS =
DOCKER_USER =

all: manager

test: generate fmt vet
	go test ./... -v -cover

# Build manager binary
manager: generate fmt vet
	go build -o gate main.go

run: generate fmt vet
	go run ./main.go

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

generate:
	go run github.com/99designs/gqlgen generate ./...

# Build the docker image
docker-build: test
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	echo $DOCKER_PASS | docker login --username $DOCKER_USER --password-stdin
	docker push ${IMG}

