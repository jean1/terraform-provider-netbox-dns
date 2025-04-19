default: fmt lint install generate

build:
	go build -v ./...

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	# TF_LOG=TRACE NETBOX_SERVER_URL=http://nb/ NETBOX_API_TOKEN=a339fbe313c1c183e7896490a9778a4981c90202 TF_ACC=1 go test -v -cover -timeout 120m ./...
	# TF_LOG=DEBUG NETBOX_SERVER_URL=http://nb/ NETBOX_API_TOKEN=a339fbe313c1c183e7896490a9778a4981c90202 TF_ACC=1 go test -v -cover -timeout 120m ./...
	NETBOX_SERVER_URL=http://nb/ NETBOX_API_TOKEN=a339fbe313c1c183e7896490a9778a4981c90202 TF_ACC=1 go test -v -cover -timeout 120m ./...

.PHONY: fmt lint test testacc build install generate

# HOSTNAME=unistra.fr
# NAMESPACE=dnum
# VERSION=0.1.0
# OS_ARCH=linux_amd64

NAME=netboxdns
BINARY=terraform-provider-${NAME}
DESTDIR = ~/.terraform.d/plugins/unistra.fr/dnum/${NAME}/0.1.0/linux_amd64

install: build
	GOSUMDB=off go build -o ${BINARY}
	mkdir -p ${DESTDIR}
	mv -v ${BINARY} ${DESTDIR}/
