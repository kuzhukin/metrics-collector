#!/bin/bash

./metricstest -test.v -test.run=^TestIteration13$ \
    -agent-binary-path=cmd/agent/agent \
    -binary-path=cmd/server/server \
    -database-dsn='postgres://postgres:12345@localhost:5431/praktikum' \
    -server-port=8080 \
    -source-path=.