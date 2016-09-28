package control

import (
	"fmt"
	"net/http"
	"strings"
)

// CURLThis takes an http.Client and http.Request and outputs the
// equivalent cURL command, to be used elsewhere.
func CURLThis(client *http.Client, req *http.Request) string {
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
