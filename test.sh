#!/usr/bin/env sh

echo "Starting test run at $(date)"

echo "Go env"
go env

echo "Running without CGO"
CGO_ENABLED=0 go test -bench . -benchmem -v -count 4 -failfast

echo "Running with CGO"
CGO_ENABLED=1 go test -bench . -benchmem -v -count 4 -failfast

echo "Running with opts"
# Add -march=native to the GOGCCFLAGS
export GOGCCFLAGS="-march=native ${GOGCCFLAGS}"
# Set GOAMD to v3 (my CPU doesn't support v4 [AVX512])
export GOAMD="v3"

CGO_ENABLED=1 go test -gccgoflags="${GOGCCFLAGS}" -bench . -benchmem -v -count 4 -failfast
