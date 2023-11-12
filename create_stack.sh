aws cloudformation create-stack \
    --stack-name TestInstanceStack \
    --template-body file://./ec2_cf.yaml \
    --capabilities CAPABILITY_IAM
aws cloudformation wait stack-create-complete --stack-name TestInstanceStack
aws cloudformation describe-stacks --stack-name TestInstanceStack
