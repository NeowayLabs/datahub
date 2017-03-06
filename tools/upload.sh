#!/bin/bash

set -o errexit
set -o nounset

echo "run: "$0" <server> <job> <dataset> <filepath>"
echo "available datasets are: (trainingset, testset.challenge, testset.result, testset.prediction)"

server=$1
job=$2
dataset=$3
file=$4

echo "uploading job: "$job" dataset: "$dataset" filepath: "$file" to : "$server

curl -v -F "$dataset=@$file" "http://$server/api/companies/jobs/$job/upload"

echo "done"
