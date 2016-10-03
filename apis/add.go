package apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/ghchinoy/atmotool/control"
)

const (
	// CMAddAPIURI is the CM endpoint for creating an API
	CMAddAPIURI = "/api/apis"
)

// NameOnlyAPI is the structure for creating an API with name only
type NameOnlyAPI struct {
	APIVersionInfo              NameValue
	AddAPIImplementationRequest CreateMechanism
}

// NameEndpointAPI is the structure for creating an API with a name and endpoint
type NameEndpointAPI struct {
	APIVersionInfo              NameValue
	AddAPIImplementationRequest ProxyImplementationRequest
}

// NameValue is a convenience for NameOnlyAPI and NameEndpointAPI
type NameValue struct {
	Name string
}

// CreateMechanism is used in NameOnlyAPI
type CreateMechanism struct {
	CreateMechanism string
}

// ProxyImplementationRequest is for creating an API with a name and endpoint
type ProxyImplementationRequest struct {
	ProxyImplementationRequest TargetEndpointURL
}

// TargetEndpointURL is a struct to hold endpoints
type TargetEndpointURL struct {
	TargetEndpointURL []string
}

// AddAPIfromExistingService publishes an existing API to the Platform
// http://docs.akana.com/cm/api/apis/m_apis_createAPI.htm
func AddAPIfromExistingService(name string, serviceID string, config control.Configuration, debug bool) error {
	if debug {
		log.Printf("Adding API - from existing service: '%s' (%s)\n", name, serviceID)
	}
	return nil
}

// AddAPIwithSpec adds in an API, given an API specification document (swagger/oai, wadl, wsdl, raml)
// http://docs.akana.com/cm/api/apis/m_apis_createAPI.htm
// this happens in two steps, first uploading the spec to the CMS staging area,
// and then adding the API, referring to the uploaded spec
func AddAPIwithSpec(name string, specpath string, config control.Configuration, debug bool) error {
	if debug {
		log.Printf("Adding API - from spec: '%s' (%s)\n", name, specpath)
	}
	// first, upload the spec doc to the CMS

	return nil
}

// AddAPINameOnly adds an API to the Platform, but with a name only - no design document
// http://docs.akana.com/cm/api/apis/m_apis_createAPI.htm
func AddAPINameOnly(name string, config control.Configuration, debug bool) error {
	if debug {
		log.Printf("Adding API - Name Only: '%s'\n", name)
	}
	nameonly := NameOnlyAPI{APIVersionInfo: NameValue{name}, AddAPIImplementationRequest: CreateMechanism{"PROXY"}}
	bytes, _ := json.Marshal(nameonly)
	if debug {
		log.Println(string(bytes))
	}
	err := postNewAPI(bytes, config, debug)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// AddNameOnlyWithEndpoint adds an API with a name and endpoint
func AddNameOnlyWithEndpoint(name string, endpoint string, config control.Configuration, debug bool) error {
	if debug {
		log.Printf("Adding API with endpoint: '%s' @ %s", name, endpoint)
	}

	var endpoints []string
	endpoints = append(endpoints, endpoint)

	targetendpoints := TargetEndpointURL{endpoints}
	proxyimpl := ProxyImplementationRequest{ProxyImplementationRequest: targetendpoints}

	nameendpoint := NameEndpointAPI{
		APIVersionInfo:              NameValue{name},
		AddAPIImplementationRequest: proxyimpl,
	}

	bytes, _ := json.Marshal(nameendpoint)
	if debug {
		log.Println(string(bytes))
	}
	err := postNewAPI(bytes, config, debug)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func postNewAPI(message []byte, config control.Configuration, debug bool) error {

	client, err := control.LoginToCM(config, debug)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	url := config.URL + CMAddAPIURI

	data := bytes.NewReader(message)

	req, err := http.NewRequest("POST", url, data)
	req.Header.Add("Content-Type", "application/vnd.soa.v81+json; charset=UTF-8")
	req.Header.Add("Accept", "application/vnd.soa.v81+json")
	req.Header.Add("Host", strings.Trim(config.URL, "https://"))
	req = control.AddCsrfHeader(req, client)
	if debug {
		log.Println("curl command:", control.CURLThis(client, req))
	}
	if debug {
		control.DebugRequestHeader(req)
	}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	if debug {
		control.DebugResponseHeader(resp)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if debug {
		fmt.Printf("%s\n", resp.Status)
		fmt.Printf("%s\n", bodyBytes)
	}

	if resp.StatusCode != 200 {
		fmt.Println(resp.Status)
	} else {
		fmt.Println("API Created ok")
	}

	return nil
}
