package apis

const (
	// APIGetInfo is the endpoint to get info about an API
	APIGetInfo = "/api/apis/%s"
	// APIGetVersionInfo is the endpoint pattern for getting information about a version, defaulting endpoint enclusion
	APIGetVersionInfo = "/api/apis/versions/%s?IncludeEndpoints=true"
	// APIVersionImplementations is the endpoint for getting implementation info
	APIVersionImplementations = "/api/apis/versions/implementations"
)
