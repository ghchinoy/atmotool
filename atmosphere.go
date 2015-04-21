package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"

	"bytes"

	"bitbucket.org/ghchinoy/atmotool/zip"
	"github.com/docopt/docopt-go"
)

// Configuration
type Configuration struct {
	Url      string `json:"url"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var config Configuration
var client *http.Client

func main() {

	usage := `SOA Software Community Manager Helper Tool.

Usage:
  atmosphere zip --prefix <prefix> --config <config> [--dir <dir>]
  atmosphere upload less <file> --config <config>
  atmosphere upload file --path <path> --config <config> <files>...
  atmosphere upload all --config <config> [--dir <dir>]
  atmosphere -h | --help
  atmosphere --version

Options:
  -h --help  Show help message and exit.
  --version  Show version and exit.
  --dir=<dir>  Directory. [default: .]
  --path=<cms_path>  CM CMS path.
`
	arguments, _ := docopt.Parse(usage, nil, true, "1.0 cirrus", false)

	// Debug for command-line args
	/*
		var keys []string
		for k := range arguments {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		// print the argument keys and values
		for _, k := range keys {
			fmt.Printf("%9s %v\n", k, arguments[k])
		}
	*/

	configLocation, _ := arguments["<config>"].(string)
	configBytes, err := ioutil.ReadFile(configLocation)
	if err != nil {
		fmt.Printf("Error opening config file %s\n", err)
		flag.Usage()
		os.Exit(1)
	}
	//var config Configuration
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		fmt.Printf("Unable to parse configuration file. %s\n", err)
		os.Exit(1)
	}

	if len(config.Password) < 1 {
		fmt.Printf("Missing or blank password.")
		os.Exit(1)
	}

	if arguments["upload"] == true {
		if arguments["less"] == true {
			lessFilePath := arguments["<file>"].(string)
			uploadLessFile(lessFilePath, config)
		} else if arguments["all"] == true {
			dir, _ := arguments["--dir"].(string)
			uploadAllHelper(dir, config)
		} else if arguments["file"] == true {
			var files []string
			for _, v := range arguments["<files>"].([]string) {
				files = append(files, v)
			}
			path, _ := arguments["--path"].(string)
			upload(files, config, path)
		}
	} else if arguments["zip"] == true {
		prefix, _ := arguments["<prefix>"].(string)
		dir, _ := arguments["--dir"].(string)
		zip.ZipPredefinedPath(prefix, dir)
	}
}

// Convenience method
func uploadLessFile(lessFilePath string, config Configuration) {
	log.Printf("Uploading Less file %s to %s\n", lessFilePath, config.Url)

	err := loginToCM()
	if err != nil {
		log.Fatalln(err)
		return
	}

	//
	// Upload
	log.Println("Uploading custom.less ...")
	extraParams := map[string]string{
		"none": "really",
	}
	lessUploadUri := config.Url + "/resources/theme/default/less?unpack=false"
	var request *http.Request
	request, err = newFileUploadRequest(lessUploadUri, extraParams, "File", lessFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		resp.Body.Close()
		log.Printf("Upload status %v", resp.StatusCode)
	}

	if resp.StatusCode == 200 {
		err = rebuildStyles("default")
		if err != nil {
			log.Fatalln(err)
		}
	}

}

func loginToCM() error {
	// Login
	log.Println("Logging in...")
	client = &http.Client{}
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	client.Jar = jar

	loginUri := config.Url + "/api/login"
	auth := Auth{config.Email, config.Password}
	buf, err := json.Marshal(auth)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	req, err := http.NewRequest("POST", loginUri, bytes.NewReader(buf))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	log.Printf("Login %s", resp.Status)

	return nil

}

func rebuildStyles(theme string) error {

	// call rebuild styles API
	// POST CM_URI/resources/branding/generatestyles
	// Form Data
	// theme: default
	log.Println("Rebuilding styles...")
	rebuildStylesUri := config.Url + "/resources/branding/generatestyles"
	resp, err := http.PostForm(rebuildStylesUri, url.Values{"theme": {theme}})

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	var results map[string]interface{}
	err = json.Unmarshal(data, &results)
	status := results["result"]
	log.Printf("Rebuild styles: %s", status)
	return nil
}

func newFileUploadRequest(uri string, params map[string]string, paramName string, path string) (*http.Request, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	//writer.SetBoundary("WebKitFormBoundaryoa9Fs4QksdotlWrl")

	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	/*
		for key, val := range params {
			_ = writer.WriteField(key, val)
		}
	*/

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("POST", uri, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	return req, nil
}

// Convenience method
func uploadAllHelper(dir string, config Configuration) {
	fmt.Printf("Uploading all in %s to %s\n", dir, config.Url)
	// upload PREFIX_resourcesThemeDefault.zip to CMS /resources/theme/default
	// upload PREFIX_contentHomeLanding.zip to CMS /content/home/landing
	// upload less file to CMS ??
	// call rebuild styles API

}

// basic upload to CMS
func upload(files []string, config Configuration, path string) {
	fmt.Printf("Uploading to %s cms location %s these: %s\n", config.Url, path, files)
	// upload FILE to CMS path PATH
	// iterate through []FILE
}
