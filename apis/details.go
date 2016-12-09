package apis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"strings"

	"github.com/ghchinoy/atmotool/cm"
	"github.com/ghchinoy/atmotool/control"
)

const (
	// APIGetInfo is the endpoint to get info about an API
	APIGetInfo = "/api/apis/%s"
	// APIGetInfoIncludeDefault is the endpoint to get info about an API and its default version
	APIGetInfoIncludeDefault = "/api/apis/%s?includeDefaultVersion=true"
	// APIGetInfoIncludeDefaultAndSettings retrieves info about the API, its default version, and settings
	APIGetInfoIncludeDefaultAndSettings = "/api/apis/%s?IncludeDefaultVersion=true&IncludeSettings=true"
	// APIGetVersionInfo is the endpoint pattern for getting information about a version, defaulting endpoint enclusion
	APIGetVersionInfo = "/api/apis/versions/%s?IncludeEndpoints=true"
	// APIVersionImplementations is the endpoint for getting implementation info
	APIVersionImplementations = "/api/apis/versions/implementations"
	// APISettings is the endpoint to retrieve only the API's settings, ref. http://docs.akana.com/cm/api/apis/m_apis_getAPISettings.htm
	APISettings = "/api/apis/%s/settings"
)

// ShowDetailsforAPIID outputs API details
func ShowDetailsforAPIID(apiID string, useVersion bool, config control.Configuration, debug bool) error {
	client, _, err := control.LoginToCM(config, debug)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	var pattern string
	if useVersion {
		pattern = APIGetVersionInfo
	} else {
		pattern = APIGetInfoIncludeDefault
	}

	url := config.URL + fmt.Sprintf(pattern, apiID)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	if debug {
		log.Println("Calling", url)
		control.DebugRequestHeader(req)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if debug {
		control.DebugResponseHeader(resp)
	}
	var apiInfo cm.ApisResponse
	err = json.Unmarshal(bodyBytes, &apiInfo)
	if err != nil {
		return err
	}
	if resp.StatusCode == 500 {
		var message string
		if strings.Contains(apiInfo.FaultMessage, "[apiversion]") {
			message = "Please provide an API ID. An API ID was expected; instead, an API Version ID was provided.\nPlease use the --ver flag."
		}
		if strings.Contains(apiInfo.FaultMessage, "[api]") {
			message = "Please provide an API Version ID. An API Version ID was expected; instead, an API ID was provided.\nPlease remove the --ver flag."
		}
		if debug {
			fmt.Printf("%s : %s\n", resp.Status, apiInfo.FaultMessage)
		}
		return errors.New(message)
	}
	if useVersion {
		var api cm.APIVersion
		err = json.Unmarshal(bodyBytes, &api)
		if err != nil {
			return err
		}
		outputAPIVersion(api)
	} else {
		var api cm.APIDetails
		err = json.Unmarshal(bodyBytes, &api)
		if err != nil {
			return err
		}
		outputAPI(api)
	}
	//fmt.Printf("%s: %s\n", resp.Status, bodyBytes)
	return nil
}

func outputAPIVersion(api cm.APIVersion) {

	fmt.Printf("API: %s (%s)\n", api.Description, api.APIID)
	fmt.Printf("Version: %s (%s)\n", api.Name, api.APIVersionID)
	for _, v := range api.Endpoints.Endpoint {
		var visibility string
		for _, t := range v.ConnectionProperties {
			if t.Name == "visibility" {
				visibility = t.Value
			}
		}
		fmt.Printf("%s (%s): %s\n", v.ImplementationCode, visibility, v.URI)
	}
}

func outputAPI(api cm.APIDetails) {

	fmt.Printf("API: %s (%s)\n", api.Description, api.APIID)
	fmt.Printf("Latest version ID: %s\n", api.LatestVersionID)
}
