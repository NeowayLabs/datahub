#!/bin/bash

set -o errexit
set -o nounset

echo "run: "$0" <server> <scriptpath>"

server=$1
script=$2

# This is all the data that the company will upload
./tools/upload.sh $server trainingset.csv ./examples/trainingset.csv
./tools/upload.sh $server testset.challenge.csv ./examples/testset.challenge.csv
#./tools/upload.sh $server testset.result.csv ./examples/testset.result.csv

# Upload the code, the statistician does that
curl -v  -F "code.r=@$script" http://$server/datahub/upload

# Run moderfocker run !!!
curl -v  http://$server/datahub/execr
