#!/bin/bash

set -o errexit
set -o nounset

echo "run: "$0" <server> <dataset> <filepath>"
echo "available datasets are: (trainingset, testset.challenge, testset.result, testset.prediction)"

server=$1
dataset=$2
file=$3

echo "uploading dataset: "$dataset" filepath: "$file" to : "$server

curl -v  -F "$dataset=@$file" http://$server/datahub/upload

echo "done"
