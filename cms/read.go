package cms

import "github.com/ghchinoy/atmotool/control"

const (
	//CMSReadFileDetailsURI is the endpoint for reading a files details
	CMSReadFileDetailsURI = "/api/dropbox/readFileDetails"
	// CMSReadURLURI is the endpoint for information for a file located at a URL
	CMSReadURLURI = "/api/dropbox/readurl"
	// CMSReadWSDLURI is the endpoint for reading a WSDL's details
	CMSReadWSDLURI = "/api/dropbox/wsdls"
)

// ReadFileDetails reads the contents of a file
// http://docs.akana.com/cm/api/dropbox/m_dropbox_readFileDetails.htm
func ReadFileDetails(config control.Configuration, debug bool) error {
	return nil
}

// ReadURL provides information about a file located at an URL
func ReadURL(config control.Configuration, debug bool) error {
	return nil
}

// ReadWSDLzip
func ReadWSDLzip(config control.Configuration, debug bool) error {
	return nil
}
