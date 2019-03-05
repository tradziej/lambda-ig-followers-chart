AWS_REGION ?= eu-central-1
S3_BUCKET ?= lambda-ig-followers-chart-deployment
STACK_NAME ?= ig-followers

TEMPLATE = template.yaml
PACKAGED_TEMPLATE = packaged.yaml

scraper: ./services/scraper/main.go
	go build -o bin/scraper ./services/scraper

drawing: ./services/drawing/main.go
	go build -o bin/drawing ./services/drawing

clean:
	rm -drf bin/

lambda:
	GOOS=linux GOARCH=amd64 $(MAKE) scraper
	GOOS=linux GOARCH=amd64 $(MAKE) drawing

build: clean lambda

create-s3-bucket:
	if ! aws s3api head-bucket --bucket $(S3_BUCKET) --region $(AWS_REGION) 2>/dev/null; then \
		aws s3 mb s3://$(S3_BUCKET) --region $(AWS_REGION) \
		; \
	fi

sam-package:
	@sam package \
    --template-file $(TEMPLATE) \
    --output-template-file $(PACKAGED_TEMPLATE) \
    --s3-bucket $(S3_BUCKET) \
		--region $(AWS_REGION)

sam-deploy:
	@sam deploy \
    --template-file $(PACKAGED_TEMPLATE) \
    --capabilities CAPABILITY_IAM \
    --stack-name $(STACK_NAME) \
		--region $(AWS_REGION)

deploy: create-s3-bucket build sam-package sam-deploy

destroy:
	@aws cloudformation delete-stack --stack-name $(STACK_NAME)

validate-template:
	@aws cloudformation validate-template \
			--template-body file://$(TEMPLATE)

describe:
	@aws cloudformation describe-stacks \
		--region $(AWS_REGION) \
		--stack-name $(STACK_NAME) \

outputs:
	@aws cloudformation describe-stacks \
		--region $(AWS_REGION) \
		--stack-name $(STACK_NAME) \
		--query 'Stacks[].Outputs'