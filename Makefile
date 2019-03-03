AWS_REGION ?= eu-central-1
S3_BUCKET ?= lambda-ig-followers-chart-deployment
STACK_NAME ?= ig-followers

create-s3-bucket:
	if ! aws s3api head-bucket --bucket $(S3_BUCKET) --region $(AWS_REGION) 2>/dev/null; then \
		aws s3 mb s3://$(S3_BUCKET) --region $(AWS_REGION) \
		; \
	fi

sam-package:
	@sam package \
    --template-file template.yaml \
    --output-template-file packaged.yaml \
    --s3-bucket $(S3_BUCKET) \
		--region $(AWS_REGION)

sam-deploy:
	@sam deploy \
    --template-file packaged.yaml \
    --capabilities CAPABILITY_IAM \
    --stack-name $(STACK_NAME) \
		--region $(AWS_REGION)

deploy: create-s3-bucket sam-package sam-deploy

destroy:
	@aws cloudformation delete-stack --stack-name $(STACK_NAME)

validate-template:
	@aws cloudformation validate-template \
			--template-body file://template.yaml

describe:
	@aws cloudformation describe-stacks \
		--region $(AWS_REGION) \
		--stack-name $(STACK_NAME) \

outputs:
	@aws cloudformation describe-stacks \
		--region $(AWS_REGION) \
		--stack-name $(STACK_NAME) \
		--query 'Stacks[].Outputs'