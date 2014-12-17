#! /bin/bash

set -xe

for i in master slave network; do go test $i/...; done
#website is not included here

echo " !!SUCCESS!! tests are done"
