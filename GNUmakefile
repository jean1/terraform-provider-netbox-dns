default: fmt lint install generate

build:
	go build -v ./...

install: build
	go install -v ./...

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

# TEST?=$$(go list ./... | grep -v 'vendor')
# HOSTNAME=unistra.fr
# NAMESPACE=dnum
# VERSION=0.1.0
# OS_ARCH=linux_amd64

# sum of packages seems to end up in https://sum.golang.org preventing some
# reinstallation of modules (like when you push --force; keep retrieving the
# same shity package somewhere ...)

NAME=netboxdns
BINARY=terraform-provider-${NAME}
#DESTDIR = /home/jean/.terraform.d/plugins/registry.terraform.io/hashicorp/netboxdns/
# /home/jean/.terraform.d/plugins/unistra.fr/dnum/${NAME}
DESTDIR = /home/jean/.terraform.d/plugins/registry.terraform.io/hashicorp/${NAME}

pbuild: 
	GOSUMDB=off go build -o ${BINARY}

#pinstall: ${BINARY}
	#mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	#mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

	# mkdir -p /usr/local/share/terraform/plugins
	# cp -v ${BINARY} /usr/local/share/terraform/plugins/

pinstall: ${BINARY}
	mkdir -p ${DESTDIR}
	cp -v ${BINARY} ${DESTDIR}/
	# to use the installed plugin: terraform init -plugin-dir='.terraform.d/plugins'

#pinit:
#	cd examples && rm -f .terraform.lock.hcl && terraform init
#
#ptest:
#	go test -i $(TEST) || exit 1
#	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4
#
#ptestacc:
#	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m
