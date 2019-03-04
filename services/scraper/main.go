package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

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

	log.Println("Followers count:", scrapeFollowersCount(body))

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
