#! /bin/bash

set -xe

for i in master network proxy slave website; do go test $i/...; done

echo " !!SUCCESS!! Tests are done. "
