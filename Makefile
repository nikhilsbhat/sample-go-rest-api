GOFMT_FILES?=$$(find . -not -path "./vendor/*" -type f -name '*.go')

default: check

prebuild:
	go get -u github.com/nikhilsbhat/go-api-sample@charts
	GOPATH=${GOPATH:-$(go env GOPATH)}
	cd GOPATH
# bin generates the releaseable binaries for config
build: prebuild check
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o sample-api

# Validates the dependencies exists.
check: fmt
	go mod vendor

# Lints all the go files if they are not.
fmt:
	gofmt -w $(GOFMT_FILES)

# Runs the application by building a binary of it.
run: build
	./config