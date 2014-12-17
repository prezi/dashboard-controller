#! /bin/bash

set -xe

go install master
go install slave
go install network
go install website

echo " !!SUCCESS!! build is ok"
