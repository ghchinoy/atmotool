package apis

import (
	"log"

	"github.com/ghchinoy/atmotool/control"
)

const (
	// GetAPIVersionsFormat golang fmt format string for http://docs.akana.com/cm/api/apis/m_apis_getAPIVersions.htm
	GetAPIVersionsFormat = "/api/apis/{APIID}/versions"
	// GetAPIVersions2Format is a 2nd golang fmt format string for http://docs.akana.com/cm/api/apis/m_apis_getAPIVersions.htm
	GetAPIVersions2Format = "/api/apis/versions/{APIVersionID}"
)

// APIVersions gets the versions of a particular API
func APIVersions(apiID string, config control.Configuration, debug bool) error {

	if debug {
		log.Println("Listing Versions of API", apiID)
	}

	return nil

}
