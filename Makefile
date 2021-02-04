NAME=spotify
BINARY=terraform-provider-${NAME}
HOSTNAME=hashicorp.com
NAMESPACE=jbonhag
OS_ARCH=darwin_amd64
VERSION=0.1

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
