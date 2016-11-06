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
	"github.com/ghchinoy/atmotool/dropbox"
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

// DLDescriptor is used to reference a previously uploaded spec doc
type APIwithSpec struct {
	DLDescriptor SDR
}

type SDR struct {
	ServiceDescriptorReference ServiceDescriptorReference
}

// ServiceDescriptorReference represents the needed information for referencing a spec doc
type ServiceDescriptorReference struct {
	ServiceName  string
	FileName     string
	DropoxFileID int `json:"DropboxFileId"`
}

// CreateAPIfromExistingService publishes an existing API to the Platform
// http://docs.akana.com/cm/api/apis/m_apis_createAPI.htm
func CreateAPIfromExistingService(name string, serviceID string, config control.Configuration, debug bool) error {
	if debug {
		log.Printf("Adding API - from existing service: '%s' (%s)\n", name, serviceID)
	}
	return nil
}

// CreateAPIwithSpec adds in an API, given an API specification document (swagger/oai, wadl, wsdl, raml)
// http://docs.akana.com/cm/api/apis/m_apis_createAPI.htm
// this happens in two steps, first uploading the spec to the CMS staging area,
// and then adding the API, referring to the uploaded spec
// This is invoked by: atmotool apis create APINAME --spec PATH_TO_SPECFILE
func CreateAPIwithSpec(name string, specpath string, config control.Configuration, debug bool) error {
	if debug {
		log.Printf("Adding API - from spec: '%s' (%s)\n", name, specpath)
	}
	// first, upload the spec doc to the CMS
	specresponse, err := dropbox.AddSpecToDropbox(config, specpath, debug)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	specref := ServiceDescriptorReference{
		ServiceName:  specresponse.ServiceDescriptorDocument[0].ServiceName[0],
		FileName:     specresponse.FileName,
		DropoxFileID: specresponse.DropboxFileID,
	}
	spec := APIwithSpec{DLDescriptor: SDR{ServiceDescriptorReference: specref}}

	bytes, _ := json.Marshal(spec)
	if debug {
		log.Println(string(bytes))
	}
	err = postNewAPI(bytes, config, debug)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	//fmt.Printf("Document ID: %v", specresponse.DropboxFileID)

	return nil
}

// CreateAPINameOnly adds an API to the Platform, but with a name only - no design document
// http://docs.akana.com/cm/api/apis/m_apis_createAPI.htm
func CreateAPINameOnly(name string, config control.Configuration, debug bool) error {
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

// CreateAPINameOnlyWithEndpoint adds an API with a name and endpoint
func CreateAPINameOnlyWithEndpoint(name string, endpoint string, config control.Configuration, debug bool) error {
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

	client, _, err := control.LoginToCM(config, debug)
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
