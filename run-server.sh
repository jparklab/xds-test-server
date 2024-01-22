#!/bin/bash
go run pkg/cmd/xds_server/xds_server.go  1>logs/server.log 2>&1
