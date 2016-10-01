package apis

// AddAPINameOnly adds an API to the Platform, but with a name only - no design document
// http://docs.akana.com/cm/api/apis/m_apis_createAPI.htm
func AddAPINameOnly() error {
	return nil
}

// AddAPIfromExistingService publishes an existing API to the Platform
// http://docs.akana.com/cm/api/apis/m_apis_createAPI.htm
func AddAPIfromExistingService() error {
	return nil
}

// AddAPIwithSpec adds in an API, given an API specification document (swagger/oai, wadl, wsdl, raml)
// http://docs.akana.com/cm/api/apis/m_apis_createAPI.htm
// this happens in two steps, first uploading the spec to the CMS staging area,
// and then adding the API, referring to the uploaded spec
func AddAPIwithSpec() error {
	return nil
}
