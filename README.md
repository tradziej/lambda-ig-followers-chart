# Lambda Instagram Followers Chart
The Serverless Application contains two functions:
* __ScraperFunction__ (scrapes instagram profile page and saves followers count in DynamoDB table)
* __DrawingFunction__ (draws the line chart which displays information about followers count over time; saves it on S3 bucket as public available object)

Each of them is scheduled to run once a day and triggered by CloudWatch Event.

## Requirements

* AWS CLI already configured with at least PowerUser permission
* [AWS SAM installed](https://docs.aws.amazon.com/lambda/latest/dg/serverless_app.html)
* [Docker installed](https://www.docker.com/community-edition)
* [Golang (v1.12)](https://golang.org)

### Local development
You need to have DynamoDB local](https://hub.docker.com/r/amazon/dynamodb-local).
```
docker run -d -p 8000:8000 amazon/dynamodb-local
```

SAM local doesn't evaluate CloudFormation conditionals so you must create DynamoDB Tables by yourself:
```
aws dynamodb create-table \
  --endpoint-url http://localhost:8000 \
  --table-name logs-table \
  --attribute-definitions \
    AttributeName=time,AttributeType=S \
    AttributeName=username,AttributeType=S \
  --key-schema \
    AttributeName=username,KeyType=HASH \
    AttributeName=time,KeyType=RANGE \
  --provisioned-throughput \
    ReadCapacityUnits=1,WriteCapacityUnits=1
```

```
aws dynamodb list-tables --endpoint-url http://localhost:8000
```

```
make build && sam local invoke ScraperFunction --log-file ./out.log --env-vars env.json --no-event
```

## Makefile
Project contains the [Makefile](Makefile) which you could use for several common tasks after customisation.

## Deployment
```
make deploy
```

### Sample graph
![Sample graph](chart.png)