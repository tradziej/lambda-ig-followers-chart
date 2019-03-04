package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/tradziej/lambda-ig-followers-chart/models"
)

var sess = session.Must(session.NewSession(&aws.Config{
	S3ForcePathStyle: aws.Bool(true),
	Region:           aws.String(os.Getenv("AWS_REGION")),
	Endpoint:         aws.String(os.Getenv("S3_ENDPOINT")),
}))

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(event events.CloudWatchEvent) error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	request, _ := http.NewRequest("GET", getRequestURL(), nil)
	updateRequestHeader(request)

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Read response data in to memory
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading HTTP body. ", err)
	}

	log.Println("Download logs from S3 to temp file")
	// Download logs from S3 to temp file
	filePath, err := downloadLogsFromS3(os.Getenv("LOGS_BUCKET"), "logs.txt")

	log.Println("Delete old logs file from S3")
	// Delete old logs file from S3
	err = deleteLogsFromS3(os.Getenv("LOGS_BUCKET"), "logs.txt")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Open downloaded file")
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	followersCount := scrapeFollowersCount(body)
	data := models.DataLog{
		Date:      time.Now(),
		Followers: followersCount,
	}
	file, _ := json.MarshalIndent(data, "", "")
	// Apend new logs
	if _, err := f.Write(file); err != nil {
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	// Upload new logs file to S3
	err = uploadLogsToS3(os.Getenv("LOGS_BUCKET"), "logs.txt", filePath)

	log.Println("Followers count:", followersCount)

	return nil
}

func downloadLogsFromS3(bucket, bucketfilePath string) (string, error) {
	downloader := s3manager.NewDownloader(sess)
	filePath := fmt.Sprintf("/tmp/temp_logs.txt")

	f, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	_, err = downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(bucketfilePath),
	})

	if err != nil {
		return "", err
	}

	defer f.Close()

	return filePath, nil
}

func deleteLogsFromS3(bucket, bucketfilePath string) error {
	svc := s3.New(sess)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(bucketfilePath),
	})
	if err != nil {
		return err
	}

	// err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
	// 	Bucket: aws.String(bucket),
	// 	Key:    aws.String(bucketfilePath),
	// })
	// if err != nil {
	// 	return err
	// }

	return nil
}

func uploadLogsToS3(bucket, bucketfilePath, filePath string) error {
	uploader := s3manager.NewUploader(sess)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(bucketfilePath),
		Body:   bufio.NewReader(file),
	})

	if err != nil {
		return err
	}

	return nil
}

func getRequestURL() string {
	endpoint := "https://www.instagram.com/{username}/"
	url := strings.Replace(endpoint, "{username}", os.Getenv("IG_USERNAME"), 1)

	return url
}

func updateRequestHeader(r *http.Request) {
	r.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 Safari/537.36")
	r.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	r.Header.Set("Referer", "https://www.google.pl")
}

func scrapeFollowersCount(body []byte) string {
	count := "0"
	// Create a regular expression to find followers count
	re := regexp.MustCompile(`(?:"userInteractionCount":")(.*?)(?:")`)
	followersCount := re.FindStringSubmatch(string(body))

	if followersCount != nil {
		count = followersCount[1]
	}

	return count
}

func main() {
	lambda.Start(Handler)
}
