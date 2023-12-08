#!/bin/bash

# Fail on error
set -e
# Parameters
instance="t4g.nano"

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
echo "Starting the upload"
# Compile for arm64 as that is the format for aws instance
GOARCH=arm64 make populator
scp -o StrictHostKeyChecking=accept-new -i ~/Downloads/graphDbIreland.pem \
    ./populator ubuntu@$ip:~/
ssh -o StrictHostKeyChecking=accept-new -i ~/Downloads/graphDbIreland.pem \
    ubuntu@$ip "./populator"
echo "Completed the upload."

# Stack teardown
echo "Triggering delete operation"
aws cloudformation delete-stack \
    --stack-name TestInstanceStack
aws cloudformation wait stack-delete-complete --stack-name TestInstanceStack
echo "Teardown complete"
