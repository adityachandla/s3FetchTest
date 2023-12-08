#!/bin/bash

# Fail on error
set -e
# Parameters
instance="c7gn.2xlarge"
parallelism="20"
runs="1000"
filename=$instance"_512Kb_"$runs"_x"$parallelism".csv"

# Stack creation
aws cloudformation create-stack \
    --stack-name TestInstanceStack \
    --template-body file://./ec2_cf.yaml \
    --parameters ParameterKey=InstanceTypeParameter,ParameterValue=$instance \
    --capabilities CAPABILITY_IAM
aws cloudformation wait stack-create-complete --stack-name TestInstanceStack
ip=$(aws cloudformation describe-stacks --stack-name TestInstanceStack | jq -r .Stacks[0].Outputs[0].OutputValue)
echo "Created instance with ip=$ip"

echo "Waiting for 10 seconds for instance initialization"
sleep 10

# Test 
echo "Starting the test"
rm bench
GOARCH=arm64 make bench
scp -o StrictHostKeyChecking=accept-new -i ~/Downloads/graphDbIreland.pem\
    ./bench ubuntu@$ip:~/
ssh -o StrictHostKeyChecking=accept-new -i ~/Downloads/graphDbIreland.pem\
    ubuntu@$ip "GOGC=500 GOMEMLIMIT=4096MiB ./bench -runs $runs -parallelism $parallelism | cat >> $filename"
scp -o StrictHostKeyChecking=accept-new -i ~/Downloads/graphDbIreland.pem\
    ubuntu@$ip:~/$filename .
echo "Completed the test and copied the file to local machine."

# Stack teardown
echo "Triggering delete operation"
aws cloudformation delete-stack \
    --stack-name TestInstanceStack
aws cloudformation wait stack-delete-complete --stack-name TestInstanceStack
echo "Teardown complete"
