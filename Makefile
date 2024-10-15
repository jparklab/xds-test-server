#
# Makefile
#
SRC=$(shell find pkg/ -type f -regex ".*.go")
PROTOSRC=$(shell find proto/ -type f -regex ".*.proto")
DIST=dist
BINDIR=$(DIST)/bin
PROTODIR=proto

all: $(BINDIR)/xds-server ${BINDIR}/config-client proto

$(BINDIR)/config-client: $(SRC)
	go build -o $@ github.com/jparklab/xds-test-server/pkg/cmd/config_client

$(BINDIR)/xds-server: $(SRC) proto
	go build -o $@ github.com/jparklab/xds-test-server/pkg/cmd/xds_server

proto: ${PROTOSRC}
	@for file in $^; do \
		protoc \
 			--go_out=. \
 			--go_opt=module=github.com/jparklab/xds-test-server \
			--go-grpc_out=. \
			--go-grpc_opt=module=github.com/jparklab/xds-test-server \
			$$file; \
	done


.phony: proto