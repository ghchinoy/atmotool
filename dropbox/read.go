package dropbox

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/net/html"

	"encoding/json"

	"strings"

	"github.com/ghchinoy/atmotool/control"
)

const (
	// DropboxReadFileDetailsURI is the endpoint for reading a files details
	DropboxReadFileDetailsURI = "/api/dropbox/readfiledetails?wrapInHTML=false&document.domain="
	// DropboxReadURLURI is the endpoint for information for a file located at a URL
	DropboxReadURLURI = "/api/dropbox/readurl"
	// DropboxReadWSDLURI is the endpoint for reading a WSDL's details
	DropboxReadWSDLURI = "/api/dropbox/wsdls"
)

// ReadFileDetailsResponse is the successful response of upload a spec file
type ReadFileDetailsResponse struct {
	FileName                  string
	FileType                  string
	DropboxFileID             int `json:"DropboxFileId"`
	ServiceDescriptorDocument []SpecDoc
}

// SpecDoc is a representation of a specification document
type SpecDoc struct {
	FileName       string
	DescriptorType string
	ServiceName    []string
}

// AddSpecToDropbox adds a document to the Platform's dropbox.
// This is the ReadFileDetails endpoint of the Dropbox Service.
// http://docs.akana.com/cm/api/dropbox/m_dropbox_readFileDetails.htm
func AddSpecToDropbox(config control.Configuration, specfilepath string, debug bool) (ReadFileDetailsResponse, error) {

	var specresponse ReadFileDetailsResponse
	if debug {
		log.Printf("Uploading %s to Platform dropbox...", specfilepath)
	}

	// log in
	client, _, err := control.LoginToCM(config, debug)
	if err != nil {
		log.Fatalln(err)
		return specresponse, err
	}

	// set url
	url := config.URL + DropboxReadFileDetailsURI
	extraParams := map[string]string{
		"none": "really",
	}

	// create a request
	req, err := constructUploadRequestForFile(url, extraParams, "FileName", specfilepath)
	req.Header.Add("Accept", "application/json, application/vnd.soa.v81+json")
	control.AddCsrfHeader(req, client)

	if debug {
		log.Println("POST to", url)
		log.Println("* URL", url)
		control.DebugRequestHeader(req)
	}

	// do the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return specresponse, err
	}

	// read response body
	// why's this bit here, since ioutil.ReadAll is used?
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	resp.Body.Close()

	b, _ := ioutil.ReadAll(body)

	// if non 200, print code
	if resp.StatusCode != 200 {
		log.Println(resp.Status)
	}
	// if debug, show headers and contents
	if debug {
		log.Println("Response")
		control.DebugResponseHeader(resp)
		log.Printf("%s", b)
	}

	specresponse, err = dealWithResponse(resp.Header.Get("Content-Type"), b)
	if err != nil {
		log.Println("Can't convert response.", err.Error())
		return specresponse, nil
	}

	return specresponse, nil
}

// parses the HTML response from adding a spec doc to the platform dropbox
func dealWithResponse(contenttype string, body []byte) (ReadFileDetailsResponse, error) {
	var rfd ReadFileDetailsResponse

	if strings.Contains(contenttype, "application/json") {
		err := json.Unmarshal(body, &rfd)
		if err != nil {
			log.Println("Can't convert response.")
			return rfd, err
		}
	} else { // try html parser
		r := bytes.NewReader(body)

		doc, err := html.Parse(r)
		if err != nil {
			log.Println("Can't parse HTML")
			return rfd, err
		}

		var bodytext string
		var f func(*html.Node, bool)
		f = func(n *html.Node, print bool) {
			/*if n.Type == html.ElementNode && n.Data == "body" {

			}*/
			print = print || (n.Type == html.ElementNode && n.Data == "body")

			if print && n.Type == html.TextNode {
				bodytext = n.Data
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c, print)
			}
		}
		f(doc, false)

		err = json.Unmarshal([]byte(bodytext), &rfd)
		if err != nil {
			log.Println("Can't convert response.")
			return rfd, err
		}

	}

	return rfd, nil

}

// returns an http.request ready to use to POST a file's contents to a URI
func constructUploadRequestForFile(uri string, params map[string]string, paramName, path string) (*http.Request, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("POST", uri, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	return req, nil

}

// ReadURL provides information about a file located at an URL
func ReadURL(config control.Configuration, debug bool) error {
	return nil
}

// ReadWSDLzip is a stub
func ReadWSDLzip(config control.Configuration, debug bool) error {
	return nil
}
