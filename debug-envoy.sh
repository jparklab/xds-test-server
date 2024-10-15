#!/bin/bash
BIN=/home/jpark/github/envoyproxy/envoy/bazel-bin/source/exe/envoy-static

# exec gdb -x gdb.cmd --args ${BIN} -c config/envoy-config.yaml \
#     --concurrency 2  \
#     --base-id 2 \

exec lldb -s lldb.cmd -- ${BIN} -c config/envoy-config.yaml \
    --concurrency 2  \
    --base-id 2
