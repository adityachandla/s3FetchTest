aws cloudformation update-stack \
    --stack-name TestInstanceStack \
    --template-body file://./ec2_cf.yaml \
    --capabilities CAPABILITY_IAM
aws cloudformation wait stack-update-complete --stack-name TestInstanceStack
