
# Image URL to use all building/pushing
IMG ?= air:latest
# Docker Hub username / password
DOCKER_PASS =
DOCKER_USER =


all:	manager

test:	fmt	vet

	@if [ -f ${HOME}/.kube/config ]; then \
		echo "";\
	else\
		mkdir ${HOME}/.kube;\
		touch ${HOME}/.kube/config;\
	fi
	go test ./... -v -cover


# Build manager binary
manager:	fmt	vet
	go build -o air main.go

run:	fmt	vet
	go run ./main.go

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Build the docker image
docker-build: test
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	echo ${DOCKER_PASS} | docker login --username ${DOCKER_USER} --password-stdin
	docker push ${IMG}

