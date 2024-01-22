#!/bin/bash
rm -f logs/envoy.log
exec ./envoy -c config/envoy-config.yaml \
    --concurrency 2  \
    --base-id 1 \
    --log-path logs/envoy.log \

# --component-log-level config:debug \

