# Set up

## Download enovy image

    RELEASE=v1.29
    IMAGE=envoyproxy/envoy:debug-${RELEASE}-latest
    docker pull ${IMAGE}
    containerID=$(docker create ${IMAGE})
    docker cp ${containerID}:/usr/local/bin/envoy ./envoy-${RELEASE}
    docker rm -v ${containerID}

    ln -s envoy-${RELEASE} envoy

## Protobuf

    go install github.com/golang/protobuf/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# How to test the server

## Start the server 

    go run pkg/cmd/xds_server/xds_server.go 

## Query the server 

    grpcurl -plaintext 127.0.0.1:50051 list

## Send commands

    go run pkg/cmd/config_client/client.go  --num-services 10
    
# Envoy

## Start Envoy

    ./run-envoy.sh

## Query config

     curl "http://127.0.0.1:9901/config_dump?include_eds=on"|jq .|less
     curl "http://127.0.0.1:9901/config_dump?resource=dynamic_active_clusters"|jq .|less
     curl "http://127.0.0.1:9901/config_dump?resource=dynamic_endpoint_configs&include_eds=on"|jq .|less
