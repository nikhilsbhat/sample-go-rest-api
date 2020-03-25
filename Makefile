GOFMT_FILES?=$$(find . -not -path "./vendor/*" -type f -name '*.go')
DOCKER_REPO="docker.io"
PROJECT_NAME=go-api-sample

default: check

# bin generates the releaseable binaries for config
build: check
	go build -ldflags="-s -w" -o ${PROJECT_NAME}

# Validates the dependencies exists.
check: fmt
	go mod vendor

# Lints all the go files if they are not.
fmt:
	gofmt -w $(GOFMT_FILES)

# Runs the application by building a binary of it.
run: build
	./${PROJECT_NAME}

docker_login:
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWD} ${DOCKER_REPO}

container: docker_login
	docker build . --tag ${DOCKER_USER}/${PROJECT_NAME}:latest

publish: container
	docker push ${DOCKER_USER}/${PROJECT_NAME}:latest