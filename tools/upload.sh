#!/bin/bash

set -o errexit
set -o nounset

server=$1
file=$2

echo "uploading local file: "$file" to : "$server

curl -v --data-binary "@$file" http://$server/datahub/upload

echo "done"
