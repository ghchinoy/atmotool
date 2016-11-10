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
	"runtime"
	"strings"

	"github.com/ghchinoy/atmotool/version"
)

// Configuration provides a simple struct to hold login info
type Configuration struct {
	URL             string `json:"url" mapstructure:"url"`
	Email           string `json:"email" mapstructure:"email"`
	Password        string `json:"password" mapstructure:"password"`
	Theme           string `json:"theme" mapstructure:"theme"`
	ConsoleUsername string `json:"consoleUsername" mapstructure:"console-username"`
	LoginDomainID   string `json:"loginDomainID"`
}

// UserInfo is the logged-in user's information
type UserInfo struct {
	UserName             string `json:"userName"`
	Status               string `json:"status"`
	AvatarURL            string `json:"avatarURL"`
	UserFDN              string `json:"userFDN"`
	LoginState           string `json:"loginState"`
	AuthTokenValidUntil  string `json:"authTokenValidUntil"`
	LoginDomainID        string `json:"loginDomainId"`
	PendingNotifications int    `json:"pendingNotifications"`
}

// InitializeConfiguration reads config (or local.conf)
func InitializeConfiguration(configLocation string, debug bool) (Configuration, error) {

	var config Configuration

	// if not specified, try ./local.conf, then try ~/.akana/local.conf
	if configLocation == "" || configLocation == "<nil>" {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Cant't get working directory,")
			return config, err
		}

		configLocation = cwd + "/local.conf"
		config, err := tryConfig(configLocation)
		if err != nil {
			fmt.Println("Couldn't find config at ./local.conf")
		} else {
			return config, nil
		}

		configLocation = userHomeDir() + "/.akana/local.conf"
		config, err = tryConfig(configLocation)
		if err != nil {
			fmt.Println("Couldn't find config at ~/.akana/local.conf")
			return config, err
		}

		return config, nil
	}

	config, err := tryConfig(configLocation)
	if err != nil {
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

// extracting to make repeatable
// for _, v := range pathstocheckarray {
//	config, err := tryConfig(v)
//  if err == nil {
//   return config
// }
//}
func tryConfig(configLocation string) (Configuration, error) {
	var config Configuration
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

	return config, nil
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// Auth is another simple struct
type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginToCM logs in to the API Platform
func LoginToCM(config Configuration, debug bool) (*http.Client, UserInfo, error) {
	var u UserInfo

	// Login
	if debug {
		log.Println("Logging in...")
	}
	client := &http.Client{}
	var err error
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
		return client, u, err
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
		return client, u, err
	}
	req, err := http.NewRequest("POST", loginURI, bytes.NewReader(buf))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return client, u, err
	}
	if resp.StatusCode != 200 || debug {
		log.Printf("Login %s", resp.Status)
	}

	// assign logindomainid to UserInfo
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if debug {
			log.Println("Can't get UserInfo")
		}
	}
	err = json.Unmarshal(bodyBytes, &u)
	if err != nil {
		if debug {
			log.Println("Cant unmarshal UserInfo")
		}
	}

	// debug
	if debug {
		DebugResponseHeader(resp)
	}

	return client, u, nil
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
		for _, h := range v {
			log.Printf("%s : %s", k, h)
		}
	}
	log.Println("<<< DEBUG <<<")

}

// DebugRequestHeader outputs headers of an http.Request struct to the log
func DebugRequestHeader(req *http.Request) {

	log.Println(">>> DEBUG >>>")
	for k, v := range req.Header {
		for _, h := range v {
			log.Printf("%s : %s", k, h)
		}
	}
	log.Println("<<< DEBUG <<<")

}
