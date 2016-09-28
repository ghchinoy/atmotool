package apis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	"github.com/ghchinoy/atmotool/cm"
	"github.com/ghchinoy/atmotool/control"
)

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

// API is a convenience structure for a CM API
type API struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	ID      string `json:"id"`
}

// APIs is a collection of API structs
type APIs []API

// Len is an implementation of sort interface for length of APIs
func (slice APIs) Len() int {
	return len(slice)
}

// Less is an implementation of sort interface for less comparison
func (slice APIs) Less(i, j int) bool {
	return slice[i].Name < slice[j].Name
}

// Swap is an implementation of the sort interface swap function
func (slice APIs) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// APIList returns a list of apis on the platform
func APIList(config control.Configuration, debug bool) error {
	//var request *http.Request
	if debug {
		log.Println("Listing APIs")
	}
	client, err := control.LoginToCM(config, debug)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	url := config.URL + CMListAPIsURI

	//client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	if debug {
		log.Println("curl command:", control.CURLThis(client, req))
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
	if debug {
		log.Printf("Found %v APIs", len(apis.Channel.Items))
	}

	var apiList APIs

	fmt.Printf("%v APIs\n", len(apis.Channel.Items))

	for _, v := range apis.Channel.Items {
		if debug {
			fmt.Printf("%s (%s)\n", v.EntityReference.Title, v.EntityReference.Guid)
		}
		apiList = append(apiList, API{Name: v.EntityReference.Title, ID: v.EntityReference.Guid})
	}

	sort.Sort(apiList)

	for _, v := range apiList {
		fmt.Printf("%-46s %-20s\n", v.ID, v.Name)
	}

	return nil
}
