#!/bin/bash

set -o errexit
set -o nounset

echo "run: "$0" <server> <scientist> <job> <scriptpath>"

server=$1
scientist=$2
job=$3
script=$3

# This is all the data that the company will upload
./tools/upload.sh $server $job trainingset.csv ./examples/trainingset.csv
./tools/upload.sh $server $job testset.challenge.csv ./examples/testset.challenge.csv
#./tools/upload.sh $server $job testset.result.csv ./examples/testset.result.csv

# Upload the code, the statistician does that
curl -v -X POST -F "code.r=@$script" "http://$server/api/scientists/$scientist/jobs/$job/upload"

# Run moderfocker run !!!
curl -v  http://$server/api/execr
