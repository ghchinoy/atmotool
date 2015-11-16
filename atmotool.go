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
	"strconv"
	"strings"

	"bitbucket.org/apihussain/atmotool/cm"
	"bitbucket.org/apihussain/atmotool/zip"

	"bytes"

	"github.com/docopt/docopt-go"
)

const (
	version     = "1.3.2"
	versionName = "cirrus"
)

// Configuration provides a simple struct to hold login info
type Configuration struct {
	Url      string `json:"url"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Auth is another simple struct
type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Api is a CM API
type Api struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Id      string `json:"id"`
}

const (
	CMLandingIndex         = "/content/home/landing/index.htm"
	CMInternationalization = "/i18n"
	CMCustomLess           = "/less/custom.less"
	CMFavicon              = "/style/images/favicon.ico"
	CMCustomLessURI        = "/resources/theme/default/less?unpack=false"
	CMListAPIsURI          = "/api/apis"
	CMListAppsURI          = "/api/apps"
	CMListPoliciesURI      = "/api/policies"
)

var (
	config Configuration
	client *http.Client
	jar    http.CookieJar
)

func main() {

	usage := `Akana Community Manager Helper Tool.

Usage:
  atmotool zip --prefix <prefix> <dir>
  atmotool upload less <file> [--config <config>]
  atmotool upload file --path <path> <files>... [--config <config>]
  atmotool download --path <path> <filename> [--config <config>]
  atmotool list apis [--config <config>]
  atmotool list apps [--config <config>]
  atmotool list policies [--config <config>]
  atmotool rebuild [<theme>] [--config <config>]
  atmotool reset [<theme>] [--config <config>]
  atmotool -h | --help
  atmotool --version

Options:
  -h --help  Show help message and exit.
  --version  Show version and exit.
  --dir=<dir>  Directory. [default: .]
  --path=<cms_path>  CM CMS path.
  --config=<config> Configuration file [default: local.conf]
`
	//   atmotool upload all --config <config> [--dir <dir>]

	arguments, _ := docopt.Parse(usage, nil, true, version+" "+versionName, false)

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
		dir, _ := arguments["<dir>"].(string)
		//zip.ZipPredefinedPath(prefix, dir)
		var fn string
		if dir == "." {
			fn = "this"
		} else {
			fn = strings.Replace(dir, ".", "", -1)
			fn = strings.Replace(fn, "/", "-", -1)
		}
		dir = strings.TrimSuffix(dir, "/")
		fn = strings.TrimSuffix(fn, "-")
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
		} else if arguments["apps"] == true {
			listApps()
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
	} else if arguments["reset"] == true {
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

		err = resetCM(theme)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		rebuildStyles(theme)
	}
}

// this function deletes an array of items in a CM
// content or resource directory
func resetCM(theme string) error {

	client, err := loginToCM()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	urls := []string{
		CMLandingIndex,
		"/resources/theme/" + theme + CMInternationalization,
		"/resources/theme/" + theme + CMCustomLess,
		"/resources/theme/" + theme + CMFavicon,
	}

	for _, url := range urls {
		urlStr := config.Url + url
		log.Println("Deleting", url)
		err := callDeleteURL(client, urlStr)
		if err != nil {
			return err
		}
	}

	return nil
}

// Used by resetCM, this is called multiple times to delete
// a specific url
func callDeleteURL(client *http.Client, urlStr string) error {

	//client := &http.Client{}
	//client.Jar = jar
	req, err := http.NewRequest("DELETE", urlStr, nil)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	//bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("Delete:", resp.Status)

	return nil
}

func listApps() error {
	log.Println("Listing Apps")

	client, err := loginToCM()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	url := config.Url + CMListAppsURI

	//client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("%s", bodyBytes)
	var apps cm.ApisResponse
	err = json.Unmarshal(bodyBytes, &apps)
	log.Printf("Found %v Apps", len(apps.Channel.Items))

	var appList []Api

	for _, v := range apps.Channel.Items {
		appList = append(appList, Api{Name: v.EntityReference.Title})
	}
	jsonBytes, err := json.Marshal(appList)
	if err != nil {
		log.Printf("Unable to marshall appList to json")
	}
	fmt.Printf("%s", jsonBytes)

	return nil
}

func listApis() error {
	//var request *http.Request
	log.Println("Listing APIs")

	client, err := loginToCM()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	url := config.Url + CMListAPIsURI

	//client := &http.Client{}
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

	client, err := loginToCM()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	url := config.Url + CMListPoliciesURI

	//client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//var policies cm.PoliciesResponse

	log.Printf("%s", bodyBytes)

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
// TODO review this - http client created, but not used?
func uploadLessFile(uploadFilePath string, config Configuration) {
	log.Printf("Uploading Less file %s to %s\n", uploadFilePath, config.Url)

	_, err := loginToCM()
	if err != nil {
		log.Fatalln(err)
		return
	}

	// Upload
	log.Println("Uploading custom.less ...")
	extraParams := map[string]string{
		"none": "really",
	}
	uploadUri := config.Url + CMCustomLessURI

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

	for k, v := range request.Header {
		log.Printf("%s : %s", k, v)
	}

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

func loginToCM() (*http.Client, error) {
	// Login
	log.Println("Logging in...")
	client = &http.Client{}
	var err error
	jar, err = cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
		return client, err
	}
	client.Jar = jar

	loginURI := config.Url + "/api/login"
	auth := Auth{config.Email, config.Password}
	buf, err := json.Marshal(auth)
	if err != nil {
		log.Fatalln(err)
		return client, err
	}
	req, err := http.NewRequest("POST", loginURI, bytes.NewReader(buf))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return client, err
	}
	if resp.StatusCode != 200 {
		log.Printf("Login %s", resp.Status)
	}

	return client, nil
}

// Call CM Rebuild Styles
func rebuildStyles(theme string) error {

	client, err := loginToCM()
	if err != nil {
		log.Fatalln(err)
		return err
	}
	// call rebuild styles API
	// POST CM_URI/resources/branding/generatestyles
	// Form Data
	// theme: default
	log.Println("Rebuilding styles...")
	rebuildStylesUri := config.Url + "/resources/branding/generatestyles"
	postdata := url.Values{}
	postdata.Set("theme", theme)

	req, _ := http.NewRequest("POST", rebuildStylesUri, bytes.NewBufferString(postdata.Encode()))
	req.Header.Add("Content-Length", strconv.Itoa(len(postdata.Encode())))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	var results map[string]interface{}
	err = json.Unmarshal(data, &results)
	//log.Println(resp.Status, results)
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

	client, err := loginToCM()
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
// TODO review this - http client created but not used?
func upload(files []string, config Configuration, path string) {
	fmt.Printf("Uploading to %s cms location %s these: %s\n", config.Url, path, files)
	// upload FILE to CMS path PATH
	// iterate through []FILE

	_, err := loginToCM()
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
