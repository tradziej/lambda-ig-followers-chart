# Lambda Instagram Followers Chart

## Requirements

* AWS CLI already configured with at least PowerUser permission
* [AWS SAM installed](https://docs.aws.amazon.com/lambda/latest/dg/serverless_app.html)
* [Docker installed](https://www.docker.com/community-edition)
* [Golang (v1.12)](https://golang.org)

### Local development
```
make build && sam local invoke ScraperFunction --log-file ./out.log --env-vars env.json --no-event
```