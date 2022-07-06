IMAGE=dreitier/silencer
EXECUTABLE=silencer

# if environment variables are defined (e.g. through GitHub Actions), use that. Otherwise, try to resolve them during build
# @see https://stackoverflow.com/a/24264930/2545275
GIT_COMMIT ?= $(shell git rev-list -1 HEAD | cut -c1-7)
GIT_TAG ?= $(shell git tag --points-at HEAD)

all: build

build:
	go mod download
	# disable CGO to use default libc
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(EXECUTABLE) -ldflags "-X main.gitCommit=$(GIT_COMMIT) -X main.gitTag=$(GIT_TAG)"
	strip $(EXECUTABLE)

clean:
	go clean
	rm -rf vendor $(EXECUTABLE)

docker-build:
	echo ${GIT_COMMIT}
	docker build --no-cache -t ${IMAGE}:${TAG} .

docker-push: docker-build
	docker tag ${IMAGE}:${TAG} ${IMAGE}:latest
	docker push ${IMAGE}:${TAG}
	docker push ${IMAGE}:latest
