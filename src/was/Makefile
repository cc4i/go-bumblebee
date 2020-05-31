
# Image URL to use all building/pushing
DOCKER_IMAGE ?= herochinese/go-bumblebee-was
DOCKER_TAG ?= $(shell git rev-parse --abbrev-ref HEAD)
# Docker Hub username / password
DOCKER_PASS =
DOCKER_USER =

all: manager

test: fmt vet
	go test ./... -v -cover

# Build manager binary
manager: fmt vet
	go build -o gate main.go

run: fmt vet
	go run ./main.go

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

#generate:
#	go run github.com/99designs/gqlgen generate ./...

# Build the docker image
docker-build:
	docker build . -t ${DOCKER_IMAGE}:${DOCKER_TAG}

# Push the docker image
docker-push:
	#echo ${DOCKER_PASS} | docker login --username ${DOCKER_USER} --password-stdin
	docker push ${DOCKER_IMAGE}:${DOCKER_TAG}

