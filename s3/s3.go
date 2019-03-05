package s3

import (
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	bucketName = os.Getenv("CHART_BUCKET")
	region     = os.Getenv("AWS_REGION")
)

type Bucket struct {
	Instance *s3.S3
}

func New() *Bucket {
	awsConfig := &aws.Config{
		Region: aws.String(region),
	}

	return &Bucket{
		Instance: s3.New(session.New(), awsConfig),
	}
}

func (b *Bucket) Upload(fileName string, buf *bytes.Buffer) (string, error) {
	if _, err := b.Instance.PutObject(&s3.PutObjectInput{
		ACL:         aws.String("public-read"),
		Body:        bytes.NewReader(buf.Bytes()),
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fileName),
		ContentType: aws.String("image/png"),
	}); err != nil {
		return "", err
	}

	return b.getObjectURL(fileName), nil
}

func (b *Bucket) getObjectURL(fileName string) string {
	return fmt.Sprintf("https://s3-%v.amazonaws.com/%v/%v", region, bucketName, fileName)
}
