package apis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/fatih/structs"
	"github.com/ghchinoy/atmotool/cm"
	"github.com/ghchinoy/atmotool/control"
)

const (
	CMLandingIndex         = "/content/home/landing/index.htm"
	CMInternationalization = "/i18n"
	CMCustomLess           = "/less/custom.less"
	CMFavicon              = "/style/images/favicon.ico"
	// CMCustomLessURI should be a template, subsitute in Configuration.Theme
	CMCustomLessURI      = "/resources/theme/default/less?unpack=false"
	CMListAPIsURI        = "/api/apis"
	CMListAppsURI        = "/api/search?sortBy=com.soa.sort.order.alphabetical&count=20&start=0&q=type:app"
	CMListPoliciesURI    = "/api/policies"
	CMListUsersURI       = "/api/search?sort=asc&sortBy=com.soa.sort.order.title_sort&Federation=false&count=20&start=0&q=type:user"
	CMListAPIVersionsURI = "/api/apis/versions"
)

// API is a convenience structure for a CM API
type API struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	Version    string `json:"version"`
	VersionID  string `json:"vid"`
	Endpoint   string `json:"endpoint"`
	Visibility string `json:"visibility"`
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

// APIListVersions outputs a list of all APIs
func APIListVersions(config control.Configuration, debug bool) error {
	if debug {
		log.Println("Listing API Versions")
	}
	client, userinfo, err := control.LoginToCM(config, debug)
	if err != nil {
		log.Fatalln(err)
		return err
	}

	url := config.URL + CMListAPIVersionsURI

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	if debug {
		log.Println("curl: ", control.CURLThis(client, req))
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

	tenantID := strings.Split(userinfo.LoginDomainID, ".")[1]

	var apiList APIs

	fmt.Printf("%v APIs\n", len(apis.Channel.Items))
	for _, v := range apis.Channel.Items {
		visibility := getVisibility(v)
		var endpoint string
		if len(v.Endpoints.Endpoint) > 0 {
			endpoint = v.Endpoints.Endpoint[0].URI
		}
		// remove that tenant suffix from API guid
		apiguid := strings.Replace(v.Guid.Value, "."+tenantID, "", -1)
		apiList = append(apiList, API{
			Version:    v.Title,
			Name:       v.EntityReferences.EntityReference[0].Title,
			ID:         apiguid,
			Endpoint:   endpoint,
			Visibility: visibility,
		})
	}

	sort.Sort(apiList)

	//pattern := "%-36s %-20s %-5s %-15s %s\n"
	pattern := fmt.Sprintf("%%-%vs %%-%vs %%-%vs %%-%vs %%-%vs\n",
		maxLengthOfField(apiList, "ID"),
		maxLengthOfField(apiList, "Name"),
		maxLengthOfField(apiList, "Version"),
		maxLengthOfField(apiList, "Visibility"),
		maxLengthOfField(apiList, "Endpoint"),
	)
	fmt.Printf(pattern, fmt.Sprintf("ID (%s)", tenantID), "Name", "Ver", "Vis", "Endpoint")

	for _, v := range apiList {
		fmt.Printf(pattern, v.ID, v.Name, v.Version, v.Visibility, v.Endpoint)
	}

	return nil
}

// APIList returns a list of apis on the platform
func APIList(config control.Configuration, debug bool) error {
	//var request *http.Request
	if debug {
		log.Println("Listing APIs")
	}
	client, userinfo, err := control.LoginToCM(config, debug)
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

	// grab tenant suffix, for removal
	if debug {
		log.Printf("LoginDomainID: %s", userinfo.LoginDomainID)
	}
	tenantID := strings.Split(userinfo.LoginDomainID, ".")[1]

	for _, v := range apis.Channel.Items {
		if debug {
			fmt.Printf("%s (%s)\n", v.EntityReference.Title, v.EntityReference.Guid)
		}
		visibility := getVisibility(v)
		// remove that tenant suffix from API guid
		apiguid := strings.Replace(v.Guid.Value, "."+tenantID, "", -1)
		versionguid := strings.Replace(v.EntityReference.Guid, "."+tenantID, "", -1)
		latestversion := strings.Replace(v.EntityReference.Title, v.Title, "", -1)
		latestversion = strings.Replace(latestversion, "(", "", -1)
		latestversion = strings.Replace(latestversion, ")", "", -1)
		latestversion = strings.TrimSpace(latestversion)
		apiList = append(apiList, API{
			Name:       v.Title,
			Version:    latestversion,
			ID:         apiguid,
			Visibility: visibility,
			VersionID:  versionguid,
		})
	}

	sort.Sort(apiList)
	pattern := fmt.Sprintf("%%-%vs %%-%vs %%-%vs %%-%vs\n",
		maxLengthOfField(apiList, "ID"),
		maxLengthOfField(apiList, "Name"),
		maxLengthOfField(apiList, "Version"),
		maxLengthOfField(apiList, "VersionID"))
	fmt.Printf(pattern, fmt.Sprintf("ID (%s)", tenantID), "Name", "Ver", "Ver ID")
	for _, v := range apiList {
		fmt.Printf(pattern, v.ID, v.Name, v.Version, v.VersionID)
	}

	return nil
}

// probably add as method on APIs struct
func maxLengthOfField(list APIs, field string) int {
	var maxlen int
	for _, v := range list {
		s := structs.Map(v)
		q := fmt.Sprintf("%v", s[field])
		if len(q) > maxlen {
			maxlen = len(q)
		}
	}
	return maxlen + 2
}

// Returns the visibility of an Item
func getVisibility(v cm.Item) string {
	var visibility string
	cats := v.Category
	for _, c := range cats {
		if c.Domain == "uddi:soa.com:visibility" {
			visibility = c.Value
		}
	}
	// Shorten Registered Users visibility
	if visibility == "com.soa.visibility.registered.users" {
		visibility = "Registered"
	}
	return visibility
}
