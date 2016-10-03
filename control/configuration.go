package control

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

	"github.com/ghchinoy/atmotool/version"
)

// Configuration provides a simple struct to hold login info
type Configuration struct {
	URL      string `json:"url"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Theme    string `json:"theme"`
}

// InitializeConfiguration reads config (or local.conf)
func InitializeConfiguration(configLocation string, debug bool) (Configuration, error) {

	var config Configuration

	if configLocation == "" || configLocation == "<nil>" {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Cant't get working directory,")
			return config, err
		}
		configLocation = cwd + "/local.conf"
	}

	configBytes, err := ioutil.ReadFile(configLocation)
	if err != nil {
		fmt.Printf("Error opening config file: %s\n", err)
		return config, err
	}
	//var config Configuration
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		fmt.Printf("Unable to parse configuration file: %s\n", err)
		return config, err
	}

	if len(config.Password) < 1 {
		fmt.Printf("Missing or blank password.")
		return config, err
	}

	if debug {
		log.Println("Config file contents:", config)
	}

	return config, nil
}

// Auth is another simple struct
type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginToCM logs in to the API Platform
func LoginToCM(config Configuration, debug bool) (*http.Client, error) {
	// Login
	if debug {
		log.Println("Logging in...")
	}
	client := &http.Client{}
	var err error
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
		return client, err
	}
	client.Jar = jar

	if debug {
		log.Println(config)
	}
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
		DebugResponseHeader(resp)
	}

	return client, nil
}

// AddCsrfHeader checks to see if cookie jar has Csrf and adds it as a header
func AddCsrfHeader(req *http.Request, client *http.Client) *http.Request {
	for _, v := range client.Jar.Cookies(req.URL) {
		if strings.HasPrefix(v.Name, "Csrf-Token") {
			req.Header.Add("X-"+v.Name, v.Value)
		}
	}
	req.Header.Add("Atmotool", version.Version())
	return req
}

// DebugResponseHeader outputs to the log the headers of an http.Response struct
func DebugResponseHeader(resp *http.Response) {

	log.Println(">>> DEBUG >>>")
	for k, v := range resp.Header {
		log.Printf("%s : %s", k, v)
	}
	log.Println("<<< DEBUG <<<")

}

func DebugRequestHeader(req *http.Request) {

	log.Println(">>> DEBUG >>>")
	for k, v := range req.Header {
		log.Printf("%s : %s", k, v)
	}
	log.Println("<<< DEBUG <<<")

}
