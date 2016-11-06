package apis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ghchinoy/atmotool/control"
)

const (
	// GetMetricsFormat is the golang fmt format string for http://docs.akana.com/cm/api/apis/m_apis_getMetrics.htm
	GetMetricsFormat = "/api/apis/versions/%s/metrics"
)

// MetricsResponse is the json object that holds metrics
type MetricsResponse struct {
	StartTime string
	EndTime   string
	Interval  []MetricCollection `json:"Interval"`
}

// MetricCollection in the MetricsResponse
type MetricCollection struct {
	StartTime string
	Metrics   []MetricNameValue `json:"Metric"`
}

// MetricNameValueCollection is a collection of MetricNameValue pairs
type MetricNameValueCollection struct {
	Metric []MetricNameValue
}

// MetricNameValue is a name:value pair
type MetricNameValue struct {
	Name  string
	Value int
}

// Metric is a convenience struct
type Metric struct {
	AvgResponseTime int
	MinResponseTime int
	MaxResponseTime int
	TotalCount      int
	SuccessCount    int
	FaultCount      int
}

// APIMetrics lists metrics of an API
func APIMetrics(apiID string, config control.Configuration, debug bool) error {
	var metrics MetricsResponse
	if debug {
		log.Println("Getting metrics for", apiID)
	}
	endpoint := fmt.Sprintf(GetMetricsFormat, apiID)
	url := fmt.Sprintf("%s%s", config.URL, endpoint)

	client, _, err := control.LoginToCM(config, debug)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if debug {
		log.Printf("%s", bodyBytes)
	}
	err = json.Unmarshal(bodyBytes, &metrics)
	if err != nil {
		return err
	}
	fmt.Println("Metrics for API ", apiID)
	format := "%-20s %-5v %-5v %-5v %-5v %-5v %-5v\n"
	fmt.Printf(format, "start", "avg", "min", "max", "tot", "succ", "fault")
	for _, v := range metrics.Interval {
		m := mapMetrics(v.Metrics)
		fmt.Printf(format,
			v.StartTime,
			m.AvgResponseTime,
			m.MinResponseTime,
			m.MaxResponseTime,
			m.TotalCount,
			m.SuccessCount,
			m.FaultCount)
	}

	return nil
}

// mapMetrics turns a metric name/value pair into a Metric object
func mapMetrics(mc []MetricNameValue) Metric {
	var m Metric
	for _, v := range mc {
		switch v.Name {
		case "avgResponseTime":
			m.AvgResponseTime = v.Value
		case "minResponseTime":
			m.MinResponseTime = v.Value
		case "maxResponseTime":
			m.MaxResponseTime = v.Value
		case "totalCount":
			m.TotalCount = v.Value
		case "successCount":
			m.SuccessCount = v.Value
		case "faultCount":
			m.FaultCount = v.Value
		}

	}
	return m
}
