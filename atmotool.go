package main

import (
	"encoding/json"
	"errors"
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

	"github.com/ghchinoy/atmotool/cm"
	"github.com/ghchinoy/atmotool/zip"

	"bytes"

	"github.com/docopt/docopt-go"
)

const (
	version     = "1.4.6"
	versionName = "cirrus"
)

// Configuration provides a simple struct to hold login info
type Configuration struct {
	URL      string `json:"url"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Theme    string `json:"theme"`
}

// Auth is another simple struct
type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// API is a CM API
type API struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	ID      string `json:"id"`
}

const (
	CMLandingIndex         = "/content/home/landing/index.htm"
	CMInternationalization = "/i18n"
	CMCustomLess           = "/less/custom.less"
	CMFavicon              = "/style/images/favicon.ico"
	// CMCustomLessURI should be a template, subsitute in Configuration.Theme
	CMCustomLessURI   = "/resources/theme/default/less?unpack=false"
	CMListAPIsURI     = "/api/apis"
	CMListAppsURI     = "/api/search?sortBy=com.soa.sort.order.alphabetical&count=20&start=0&q=type:app"
	CMListPoliciesURI = "/api/policies"
	CMListUsersURI    = "/api/search?sort=asc&sortBy=com.soa.sort.order.title_sort&Federation=false&count=20&start=0&q=type:user"
)

var (
	config Configuration
	client *http.Client
	jar    http.CookieJar
	debug  bool
)

func main() {

	usage := `Akana Community Manager Command-Line Interface.

Usage:
  atmotool zip --prefix <prefix> <dir>
  atmotool upload less <file> [--config <config>] [--debug]
  atmotool upload file --path <path> <files>... [--config <config>] [--debug]
  atmotool download --path <path> <filename> [--config <config>] [--debug]
  atmotool list apis [--config <config>] [--debug]
  atmotool list topapis [--config <config>] [--debug]
  atmotool list apps [--config <config>] [--debug]
  atmotool list users [--config <config>] [--debug]
  atmotool list policies [--config <config>] [--debug]
  atmotool list cms [<path>] [--config <config>] [--debug]
  atmotool rebuild [<theme>] [--config <config>] [--debug]
  atmotool reset [<theme>] [--config <config>] [--debug]
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

	if arguments["--debug"] == true {
		debug = true
		log.Println("Debug output requested.")
	}

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

		theme := config.Theme

		if theme == "" {
			theme = "default"
		}
		// override config from cmdline
		theme, _ = arguments["<theme>"].(string)

		log.Println("Rebuilding styles for theme:", theme)

		err = rebuildStyles(theme)
		if err != nil {
			log.Println(err)
		}
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
		} else if arguments["users"] == true {
			listUsers()
		} else if arguments["topapis"] == true {
			listTopApis()
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

		// config file
		theme := config.Theme
		if theme == "" {
			theme = "default"
		}
		// override from cmdline
		theme, _ = arguments["<theme>"].(string)

		err = resetCM(theme)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		rebuildStyles(theme)
	}
}

// resetCM deletes an array of items in a CM
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
		"/resources/theme/" + theme + "/SOA",
		"/resources/theme/" + theme + "/less",
		"/resources/theme/" + theme + "/style",
	}

	for _, url := range urls {
		urlStr := config.URL + url
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
	addCsrfHeader(req, client)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	//bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("Delete:", resp.Status)

	return nil
}

// curlThis takes an http.Client and http.Request and outputs the
// equivalent cURL command, to be used elsewhere.
func curlThis(client *http.Client, req *http.Request) string {
	curl := "curl -v"
	for _, v := range client.Jar.Cookies(req.URL) {
		if strings.HasPrefix(v.Name, "Csrf-Token") {
			curl += fmt.Sprintf(" -H \"X-%s: %s\"", v.Name, v.Value)
		} else {
			curl += fmt.Sprintf(" --cookie \"%s\"", v)
		}
	}
	for k, v := range req.Header {
		for _, hv := range v {
			curl += fmt.Sprintf(" -H \"%s:%v\"", k, hv)
		}
	}
	curl += fmt.Sprintf(" %s", req.URL)

	return curl
}

// May not work with 8.0, /api/apps removed?
func listApps() error {
	log.Println("Listing Apps")

	client, err := loginToCM()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	url := config.URL + CMListAppsURI

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
	var apps cm.ApisResponse
	err = json.Unmarshal(bodyBytes, &apps)
	log.Printf("Found %v Apps", len(apps.Channel.Items))

	var appList []API

	for _, v := range apps.Channel.Items {
		appList = append(appList, API{Name: v.Title})
	}
	jsonBytes, err := json.Marshal(appList)
	if err != nil {
		log.Printf("Unable to marshall appList to json")
	}
	fmt.Printf("%s", jsonBytes)

	return nil
}

func listUsers() error {
	log.Println("Listing Users")

	client, err := loginToCM()
	if err != nil {
		return err
	}

	url := config.URL + CMListUsersURI

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	log.Println("curl command:", curlThis(client, req))
	resp, err := client.Do(req)
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("%s", bodyBytes)

	return nil
}

func listTopApis() error {
	log.Println("Listing Top APIs")

	client, err := loginToCM()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	url := config.URL + "/api/businesses/tenantbusiness.enterpriseapi/metrics?TimeInterval=15m&Duration=all&Environment=All&ReportType=business.top10.apis"

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	if debug {
		log.Println("curl command: ", curlThis(client, req))
	}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("%s", bodyBytes)
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

	url := config.URL + CMListAPIsURI

	//client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	if debug {
		log.Println("curl command:", curlThis(client, req))
	}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if debug {
		fmt.Printf("%s", bodyBytes)
	}
	var apis cm.ApisResponse
	err = json.Unmarshal(bodyBytes, &apis)
	log.Printf("Found %v APIs", len(apis.Channel.Items))

	var apiList []API

	for _, v := range apis.Channel.Items {
		if debug {
			fmt.Printf("%s (%s)\n", v.EntityReference.Title, v.EntityReference.Guid)
		}
		apiList = append(apiList, API{Name: v.EntityReference.Title, ID: v.EntityReference.Guid})
	}

	jsonBytes, err := json.Marshal(apiList)
	if err != nil {
		log.Printf("Unable to marshall apilist to json")
	}
	fmt.Printf("%s", jsonBytes)

	return nil
}

// incomplete - list raw json of policies
// should show a more human readable output
func listPolicies() error {
	log.Println("Listing Policies")

	client, err := loginToCM()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	policyTypes := []string{"Operational Policy", "Denial of Service", "Compliance Policy", "Service Level Policy"}

	for _, policyType := range policyTypes {
		log.Printf("%s\n", policyType)
		url := config.URL + CMListPoliciesURI + "?Type=" + url.QueryEscape(policyType)
		//log.Printf("* %s\n", url)

		//client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Add("Accept", "application/json")
		resp, err := client.Do(req)

		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var policies cm.ApisResponse
		err = json.Unmarshal(bodyBytes, &policies)
		log.Println("Found", len(policies.Channel.Items), " policies.")

		if len(policies.Channel.Items) > 1 {
			//log.Printf("%s", bodyBytes)
			for _, v := range policies.Channel.Items {
				fmt.Printf("\"%s\" (%s)\n", v.Title, v.Guid.Value)
			}
		}

	}

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

	if debug {
		log.Println("Config file contents:", config)
	}

	return nil
}

// Convenience method
// TODO review this - http client created, but not used?
func uploadLessFile(uploadFilePath string, config Configuration) {
	log.Printf("Uploading Less file %s to %s\n", uploadFilePath, config.URL)

	client, err := loginToCM()
	if err != nil {
		log.Fatalln(err)
		return
	}

	// Upload
	log.Println("Uploading custom.less ...")
	extraParams := map[string]string{
		"none": "really",
	}
	uploadURI := config.URL + CMCustomLessURI
	if config.Theme != "" {
		uploadURI = config.URL + "/resources/theme/" + config.Theme + "/less?unpack=false"
	}

	statusCode, err := uploadFile(client, uploadFilePath, extraParams, uploadURI)
	if err != nil {
		log.Fatalf("Issues. %v : %s", statusCode, err)
	}

	log.Printf("Upload status %v", statusCode)

	if statusCode == 200 {
		err = rebuildStyles(config.Theme)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func uploadFile(client *http.Client, uploadFilePath string, extras map[string]string, uploadURI string) (int, error) {
	var uploadStatus int

	//var request *http.Request
	request, err := newFileUploadRequest(uploadURI, extras, "File", uploadFilePath)
	if err != nil {
		log.Fatalln(err)
		return uploadStatus, err
	}
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	addCsrfHeader(request, client)

	// debug
	if debug {
		log.Println("* URL", uploadURI)
		log.Println("* Upload Path", uploadFilePath)
		for k, v := range request.Header {
			log.Printf("* %s: %s", k, v)
		}
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return uploadStatus, err
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	resp.Body.Close()

	uploadStatus = resp.StatusCode

	if uploadStatus != 200 {
		b, _ := ioutil.ReadAll(body)
		log.Println("* uploadFile", string(b))
	}

	//log.Printf("Upload status %v", resp.StatusCode)

	return uploadStatus, nil
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
		// for any extra params, map of string keys and vals
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
	fmt.Printf("Uploading all in %s to %s\n", dir, config.URL)
}

func loginToCM() (*http.Client, error) {
	// Login
	if debug {
		log.Println("Logging in...")
	}
	client = &http.Client{}
	var err error
	jar, err = cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
		return client, err
	}
	client.Jar = jar

	loginURI := config.URL + "/api/login"
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

	// debug
	if debug {
		log.Println(">>> DEBUG >>>")
		for k, v := range resp.Header {
			log.Printf("%s : %s", k, v)
		}
		log.Println("<<< DEBUG <<<")
	}

	return client, nil
}

// checks to see if cookie jar has Csrf and adds it as a header
func addCsrfHeader(req *http.Request, client *http.Client) *http.Request {
	for _, v := range client.Jar.Cookies(req.URL) {
		if strings.HasPrefix(v.Name, "Csrf-Token") {
			req.Header.Add("X-"+v.Name, v.Value)
		}
	}
	req.Header.Add("Atmotool", version)
	return req
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
	if len(theme) == 0 {
		theme = "default"
	}
	log.Printf("Rebuilding styles for theme %s ...\n", theme)
	rebuildStylesURI := config.URL + "/resources/branding/generatestyles"
	postdata := url.Values{}
	postdata.Set("theme", theme)

	req, _ := http.NewRequest("POST", rebuildStylesURI, bytes.NewBufferString(postdata.Encode()))
	req.Header.Add("Content-Length", strconv.Itoa(len(postdata.Encode())))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	addCsrfHeader(req, client)
	resp, err := client.Do(req)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	if resp.StatusCode != 200 {
		if debug {
			log.Println(resp.StatusCode, resp.Status)
		}
		return errors.New("Unable to parse less file. " + resp.Status + " when calling " + rebuildStylesURI)
	}

	var results map[string]interface{}
	err = json.Unmarshal(data, &results)
	status := results["result"]
	log.Printf("Rebuild styles: %s", status)
	return nil
}

// Download a CMS path to file
func download(path string, outputFilename string) {
	fmt.Printf("Downloading CMS path %s to file %s\n", path, outputFilename)

	client, err := loginToCM()
	if err != nil {
		log.Fatalln(err)
		return
	}

	downloadUri := config.URL + path + "?download=true&Zip=true"

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
	fmt.Printf("Uploading to %s cms location %s these: %s\n", config.URL, path, files)
	// upload FILE to CMS path PATH
	// iterate through []FILE

	client, err := loginToCM()
	if err != nil {
		log.Fatalln(err)
		return
	}

	uploadURI := config.URL + path

	extraParams := map[string]string{
		"none": "really",
	}

	for _, v := range files {
		log.Printf("Uploading %s ...\n", v)
		if strings.HasSuffix(v, ".zip") {
			uploadURI += "?unpack=true"
		}
		//log.Println(uploadURI)
		statusCode, err := uploadFile(client, v, extraParams, uploadURI)
		if err != nil {
			log.Fatalf("Issues. %v : %s", statusCode, err)
		}
		log.Printf("Upload status %v", statusCode)
	}
}
