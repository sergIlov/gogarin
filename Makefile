APP?=gogarin
SERVICE?=space_center
RELEASE?=0.1.0
GOOS?=linux
GOARCH?=amd64
DISTDIR?=dist

COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

ERRORS_ONLY?=""

.PHONY: lint
lint: prepare_metalinter
	gometalinter \
		--enable=megacheck \
		--enable=gofmt \
		--enable=goimports \
		--enable=lll --line-length=120 \
		--enable=misspell \
		--enable=unparam \
		--tests \
		$(shell test -n "${ERRORS_ONLY}" && echo --errors) \
		--vendor ./...

.PHONY: build
build: clean
	GOOS=${GOOS} GOARCH=${GOARCH} go build \
		-ldflags "-X main.version=${RELEASE} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}" \
		-o ${DISTDIR}/${APP}/${SERVICE}-${RELEASE}-${GOOS}-${GOARCH} \
		"./cmd/${SERVICE}"

.PHONY: clean
clean:
	@rm -f ${DISTDIR}/${APP}/${SERVICE}-${RELEASE}-${GOOS}-${GOARCH}

.PHONY: vendor
vendor: prepare_dep
	dep ensure -vendor-only

HAS_DEP := $(shell command -v dep;)
HAS_METALINTER := $(shell command -v gometalinter;)

.PHONY: prepare_dep
prepare_dep:
ifndef HAS_DEP
	go get -u -v -d github.com/golang/dep/cmd/dep && \
	go install -v github.com/golang/dep/cmd/dep
endif

.PHONY: prepare_metalinter
prepare_metalinter:
ifndef HAS_METALINTER
	go get -u -v -d github.com/alecthomas/gometalinter && \
	go install -v github.com/alecthomas/gometalinter && \
	gometalinter --install --update
endif