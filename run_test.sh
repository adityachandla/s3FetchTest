#!/bin/bash

# Fail on error
set -e
# Parameters
instance="t3.nano"
size="512"
runs="1000"
filename=$instance"_"$size"b_"$runs".csv"

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
make bench
scp -o StrictHostKeyChecking=accept-new -i ~/Downloads/graphDb.pem ./bench ubuntu@$ip:~/
ssh -o StrictHostKeyChecking=accept-new -i ~/Downloads/graphDb.pem ubuntu@$ip "./bench -runs $runs -size $size | cat >> $filename"
scp -o StrictHostKeyChecking=accept-new -i ~/Downloads/graphDb.pem ubuntu@$ip:~/$filename .
echo "Completed the test and copied the file to local machine."

# Stack teardown
echo "Triggering delete operation"
aws cloudformation delete-stack \
    --stack-name TestInstanceStack
aws cloudformation wait stack-delete-complete --stack-name TestInstanceStack
echo "Teardown complete"
