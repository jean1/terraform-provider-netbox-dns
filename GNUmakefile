# Plugin name
NAME=netboxdns
# Test data
URL=http://nb/
TOKEN=a339fbe313c1c183e7896490a9778a4981c90202
# Plugin path
HOSTNAME=unistra.fr
NAMESPACE=dnum
VERSION=0.1.0
OS_ARCH=linux_amd64

BINARY=terraform-provider-${NAME}
DESTDIR = ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

.PHONY: fmt lint test testacc build install
default: fmt lint install
build:
	go build -v ./...
lint:
	golangci-lint run
fmt:
	gofmt -s -w -e .
test:
	go test -v -cover -timeout=120s -parallel=10 ./...
testacc:
	NETBOX_SERVER_URL=${URL} NETBOX_API_TOKEN=${TOKEN} TF_ACC=1 go test -v -cover -timeout 120m ./...
install:
	GOSUMDB=off go build -o ${BINARY}
	mkdir -p ${DESTDIR}
	mv -v ${BINARY} ${DESTDIR}/
