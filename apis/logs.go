package apis

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ghchinoy/atmotool/control"
)

const (
	// CMExportUsageLogsFormat is a string format for the endpoint http://docs.akana.com/cm/api/apis/m_apis_exportUsageLogs.htm
	CMExportUsageLogsFormat = "/api/apis/versions/%s/txlogs/export"
)

// APILogs lists logs for an API
func APILogs(apiID string, config control.Configuration, debug bool) error {
	endpoint := fmt.Sprintf(CMExportUsageLogsFormat, apiID)
	url := fmt.Sprintf("%s%s", config.URL, endpoint)
	if debug {
		log.Println("Listing logs for API", apiID)
		log.Printf("Endpoint: %s", url)
	}

	client, err := control.LoginToCM(config, debug)
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
		log.Println(resp.Header.Get("Content-Type"))
		log.Printf("%s", bodyBytes)
	}

	fmt.Printf("%s\n", bodyBytes)

	return nil
}
