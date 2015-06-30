package main

import (
	"encoding/json"
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
	"strings"

	"bitbucket.org/apihussain/atmotool/cm"
	"bitbucket.org/apihussain/atmotool/zip"

	"bytes"

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

type Api struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Id      string `json:"id"`
}

const (
	CUSTOM_LESS_URI = "/resources/theme/default/less?unpack=false"
	LIST_APIS_URI   = "/api/apis"
)

var (
	config Configuration
	client *http.Client
)

func main() {

	usage := `Akana Community Manager Helper Tool.

Usage:
  atmosphere zip --prefix <prefix> [--dir <dir>]
  atmosphere upload less <file> [--config <config>]
  atmosphere upload file --path <path> <files>... [--config <config>]
  atmosphere download --path <path> <filename> [--config <config>]
  atmosphere list apis [--config <config>]
  atmosphere list policies [--config <config>]
  atmosphere rebuild [<theme>] [--config <config>]
  atmosphere -h | --help
  atmosphere --version

Options:
  -h --help  Show help message and exit.
  --version  Show version and exit.
  --dir=<dir>  Directory. [default: .]
  --path=<cms_path>  CM CMS path.
  --config=<config> Configuration file [default: local.conf]
`
	//   atmosphere upload all --config <config> [--dir <dir>]

	arguments, _ := docopt.Parse(usage, nil, true, "1.1.0 cirrus", false)

	// Debug for command-line args
	/*
		var keys []string
		for k := range arguments {
			keys = append(keys, k)
		}
		//sort.Strings(keys)
		// print the argument keys and values
		for _, k := range keys {
			fmt.Printf("%9s %v\n", k, arguments[k])
		}
	*/

	// convert to switch?

	if arguments["upload"] == true {
		configLocation, _ := arguments["--config"].(string)
		err := initializeConfiguration(configLocation)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if arguments["less"] == true {
			uploadFilePath := arguments["<file>"].(string)
			uploadLessFile(uploadFilePath, config)
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
		//zip.ZipPredefinedPath(prefix, dir)
		var fn string
		if dir == "." {
			fn = "this"
		} else {
			fn = strings.Replace(dir, ".", "", -1)
			fn = strings.Replace(fn, "/", "-", -1)
		}
		fn = prefix + "_" + fn + ".zip"
		fmt.Printf("Zipping %s as %s...\n", dir, fn)

		zip.ZipFolder(dir, fn)
	} else if arguments["rebuild"] == true {

		configLocation, _ := arguments["--config"].(string)
		err := initializeConfiguration(configLocation)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		theme, _ := arguments["<theme>"].(string)
		if theme == "" {
			theme = "default"
		}

		log.Println("Rebuilding styles for theme:", theme)

		rebuildStyles(theme)
	} else if arguments["list"] == true {
		// List policies
		// List APIs

		configLocation, _ := arguments["--config"].(string)
		err := initializeConfiguration(configLocation)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if arguments["policies"] == true {
			listPolicies()
		} else if arguments["apis"] == true {
			listApis()
		}
	} else if arguments["download"] == true {
		// Download path as filename.zip
		configLocation, _ := arguments["--config"].(string)
		err := initializeConfiguration(configLocation)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		path, _ := arguments["--path"].(string)
		outputFilename := arguments["<filename>"].(string)
		if !strings.HasSuffix(outputFilename, ".zip") {
			outputFilename += ".zip"
		}
		download(path, outputFilename)
	}
}

func listApis() error {
	//var request *http.Request
	log.Println("Listing APIs")

	err := loginToCM()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	url := config.Url + LIST_APIS_URI

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var apis cm.ApisResponse
	err = json.Unmarshal(bodyBytes, &apis)
	log.Printf("Found %v APIs", len(apis.Channel.Items))

	var apiList []Api

	for _, v := range apis.Channel.Items {
		//fmt.Printf("%s (%s)\n", v.EntityReference.Title, v.EntityReference.Guid)
		apiList = append(apiList, Api{Name: v.EntityReference.Title})
	}

	jsonBytes, err := json.Marshal(apiList)
	if err != nil {
		log.Printf("Unable to marshall apilist to json")
	}
	fmt.Printf("%s", jsonBytes)

	return nil
}

func listPolicies() error {
	log.Println("Listing Policies")
	return nil
}

// Reads config (or local.conf)
func initializeConfiguration(configLocation string) error {

	if configLocation == "" || configLocation == "<nil>" {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Cant't get working directory,")
			return err
		}
		configLocation = cwd + "/local.conf"
	}

	configBytes, err := ioutil.ReadFile(configLocation)
	if err != nil {
		fmt.Printf("Error opening config file: %s\n", err)
		return err
	}
	//var config Configuration
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		fmt.Printf("Unable to parse configuration file: %s\n", err)
		return err
	}

	if len(config.Password) < 1 {
		fmt.Printf("Missing or blank password.")
		return err
	}

	return nil
}

// Convenience method
func uploadLessFile(uploadFilePath string, config Configuration) {
	log.Printf("Uploading Less file %s to %s\n", uploadFilePath, config.Url)

	err := loginToCM()
	if err != nil {
		log.Fatalln(err)
		return
	}

	// Upload
	log.Println("Uploading custom.less ...")
	extraParams := map[string]string{
		"none": "really",
	}
	uploadUri := config.Url + CUSTOM_LESS_URI

	statusCode, err := uploadFile(uploadFilePath, extraParams, uploadUri)
	if err != nil {
		log.Fatalf("Issues. %v : %s", statusCode, err)
	}

	log.Printf("Upload status %v", statusCode)

	if statusCode == 200 {
		err = rebuildStyles("default")
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func uploadFile(uploadFilePath string, extras map[string]string, uploadUri string) (int, error) {
	var uploadStatus int

	var request *http.Request
	request, err := newFileUploadRequest(uploadUri, extras, "File", uploadFilePath)
	if err != nil {
		log.Fatalln(err)
		return uploadStatus, err
	}
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	// debug
	/*
		for k, v := range request.Header {
			log.Printf("%s : %s", k, v)
		}
	*/

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return uploadStatus, err
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		resp.Body.Close()

		uploadStatus = resp.StatusCode

		if uploadStatus != 200 {
			b, _ := ioutil.ReadAll(body)
			log.Println(string(b))
		}

		//log.Printf("Upload status %v", resp.StatusCode)
	}
	return uploadStatus, nil
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
	if resp.StatusCode != 200 {
		log.Printf("Login %s", resp.Status)
	}

	return nil
}

// Call CM Rebuild Styles
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
// TODO
// upload PREFIX_resourcesThemeDefault.zip to CMS /resources/theme/default
// upload PREFIX_contentHomeLanding.zip to CMS /content/home/landing
// upload less file to CMS ??
// call rebuild styles API
func uploadAllHelper(dir string, config Configuration) {
	fmt.Printf("Uploading all in %s to %s\n", dir, config.Url)
}

// Download a CMS path to file
func download(path string, outputFilename string) {
	fmt.Printf("Downloading CMS path %s to file %s\n", path, outputFilename)

	err := loginToCM()
	if err != nil {
		log.Fatalln(err)
		return
	}

	downloadUri := config.Url + path + "?download=true&Zip=true"

	file, err := os.Create(outputFilename)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer file.Close()

	/*
		check := http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
		}
	*/
	resp, err := client.Get(downloadUri)
	if err != nil {
		log.Fatalln(err)
		return
	}
	if resp.StatusCode != 200 {
		log.Fatalln(resp.StatusCode, "Unauthorized access to", downloadUri)
		return
	}
	defer resp.Body.Close()
	log.Println(resp.Status)

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Printf("%s with %v bytes downloaded.", outputFilename, size)
}

// basic upload to CMS
func upload(files []string, config Configuration, path string) {
	fmt.Printf("Uploading to %s cms location %s these: %s\n", config.Url, path, files)
	// upload FILE to CMS path PATH
	// iterate through []FILE

	err := loginToCM()
	if err != nil {
		log.Fatalln(err)
		return
	}

	uploadUri := config.Url + path

	extraParams := map[string]string{
		"none": "really",
	}

	for _, v := range files {
		log.Printf("Uploading %s ...\n", v)
		if strings.HasSuffix(v, ".zip") {
			uploadUri += "?unpack=true"
		}
		//log.Println(uploadUri)
		statusCode, err := uploadFile(v, extraParams, uploadUri)
		if err != nil {
			log.Fatalf("Issues. %v : %s", statusCode, err)
		}
		log.Printf("Upload status %v", statusCode)
	}
}
