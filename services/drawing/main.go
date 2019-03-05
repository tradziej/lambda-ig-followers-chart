package main

import (
	"bytes"
	"log"
	"strconv"
	"time"

	"github.com/tradziej/lambda-ig-followers-chart/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/tradziej/lambda-ig-followers-chart/db"
	"github.com/tradziej/lambda-ig-followers-chart/s3"
	"github.com/wcharczuk/go-chart"
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(event events.CloudWatchEvent) error {
	tbl := db.New()
	items, err := tbl.GetItems()
	if err != nil {
		log.Fatal(err)
	}

	chartBuf, err := drawChart(&items)
	if err != nil {
		log.Fatal("Error on drawing chart ", err)
	}

	bucket := s3.New()

	chartURL, err := bucket.Upload("chart.png", chartBuf)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Chart URL:", chartURL)

	return err
}

func drawChart(data *[]models.DataLog) (*bytes.Buffer, error) {
	var xval []time.Time
	var yval []float64
	var min, max float64

	for _, v := range *data {
		xval = append(xval, v.Time)
		followers, _ := strconv.ParseFloat(v.Followers, 64)
		yval = append(yval, followers)

		if followers > max {
			max = followers
		}
		if followers < min {
			min = followers
		}
	}

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
			Range: &chart.ContinuousRange{
				Min: min,
				Max: max,
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: xval,
				YValues: yval,
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
