package main

import (
	"bytes"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tradziej/lambda-ig-followers-chart/s3"
	"github.com/wcharczuk/go-chart"
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(event events.CloudWatchEvent) error {
	chartBuf, err := drawChart()
	if err != nil {
		log.Fatal("Error on drawing chart", err)
	}

	bucket := s3.New()

	chartURL, err := bucket.Upload("chart.png", chartBuf)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Chart URL:", chartURL)

	return err
}

func drawChart() (*bytes.Buffer, error) {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:      "Date",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Name:      "Followers count",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: []time.Time{
					time.Now().AddDate(0, 0, -6),
					time.Now().AddDate(0, 0, -5),
					time.Now().AddDate(0, 0, -4),
					time.Now().AddDate(0, 0, -3),
					time.Now().AddDate(0, 0, -2),
					time.Now().AddDate(0, 0, -1),
					time.Now(),
				},
				YValues: []float64{10000, 10501, 11123, 11767, 12000, 12500, 13021},
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)

	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func main() {
	lambda.Start(Handler)
}
