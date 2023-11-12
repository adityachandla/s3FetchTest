aws cloudformation delete-stack \
    --stack-name TestInstanceStack
aws cloudformation wait stack-delete-complete --stack-name TestInstanceStack
