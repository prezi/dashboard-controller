#! /bin/bash

set -xe

. $(dirname $0)/common

echo $(go version)
export GOPATH=$root
echo $GOPATH
cd "${root}"

go get github.com/stretchr/testify
go get github.com/gorilla/securecookie
go get github.com/gorilla/mux

echo "go get finished, dependencies installed"

./scripts/run_all_tests.sh
./scripts/build.sh
