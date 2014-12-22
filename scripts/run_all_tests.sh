#! /bin/bash

set -xe

for i in master slave network website; do go test $i/...; done

echo " !!SUCCESS!! Tests are done. "
