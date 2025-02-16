package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := 30 * time.Second

	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    "http://localhost:8080/api/sendCoin",
		Body:   []byte(`{"toUser": "nglmq1", "amount": 1}`),
		Header: http.Header{"Authorization": []string{"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk3Mzc5MjUsIlVzZXJJRCI6Im5nbG1xIn0.ePzUhW8G8ucmqd_6N6wYG3uLlrBCSaaljDKqFBbv-uM"}}, // Замените на ваш реальный JWT
	})

	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()

	reportFile, err := os.Create("vegeta_report.json")
	if err != nil {
		fmt.Printf("Error creating report file: %v\n", err)
		return
	}
	defer reportFile.Close()

	report := vegeta.NewJSONReporter(&metrics)
	err = report.Report(reportFile)
	if err != nil {
		fmt.Printf("Error generating report: %v\n", err)
		return
	}

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	fmt.Printf("Mean latency: %s\n", metrics.Latencies.Mean)
	fmt.Printf("Max latency: %s\n", metrics.Latencies.Max)
	fmt.Printf("Requests per second: %f\n", metrics.Rate)

	if metrics.Latencies.P99.Seconds()*1000 > 50 {
		fmt.Println("WARNING: 99th percentile latency exceeds 50ms")
	} else {
		fmt.Println("99th percentile latency is within the acceptable range.")
	}

	successRate := float64(metrics.Success) * 100 / float64(metrics.Requests)
	if successRate < 99.99 {
		fmt.Println("WARNING: Success rate is below 99.99%")
	} else {
		fmt.Println("Success rate is within the acceptable range.")
	}
}
