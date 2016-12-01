package cm

import "time"

// ApisResponse is the main struct for the RSS feed
type ApisResponse struct {
	Channel      Channel `json:"channel"`
	FaultCode    string  `json:"faultcode"`
	FaultMessage string  `json:"faultstring"`
}

// Channel is the container for items
type Channel struct {
	Title string `json:"title"`
	Items []Item `json:"item"`
}

// Item is the generic container
type Item struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Category    []ValueDomain `json:"category"`
	Guid        Guid          `json:"guid"`
	PubDate     string        `json:"pubDate"`
	// Note, the XML field tag notation does not work for JSON
	//EntityReferences []EntityReference `json:"EntityReferences.EntityReference"`
	EntityReferences struct {
		EntityReference []EntityReference
	}
	EntityReference EntityReference `json:"EntityReference"`
	ImageUrl        string          `json:"Image.Url"`
	// App
	Connections int
	Followers   int
	Rating      float32
	// User
	LastLogin     int
	Email         string
	UserName      string
	ApisCount     int
	AppsCount     int
	PostsCount    int
	CommentsCount int
	GroupsCount   int
	Domain        string
	//Endpoints     []Endpoint `json:"Endpoints.Endpoint"`
	Endpoints struct {
		Endpoint []Endpoint
	}
}

// ValueDomain is a key - value pair
type ValueDomain struct {
	Value  string `json:"value"`
	Domain string `json:"domain"`
}

// Guid is a guid string
type Guid struct {
	Value string `json:"value"`
}

// EntityReference is a reference to another entity
type EntityReference struct {
	Title    string
	Guid     string `json:"Guid"`
	Category ValueDomain
}

// Endpoint is a structure of an endpoint
type Endpoint struct {
	BindingQName         string
	BindingType          string
	CName                string
	Category             string
	ConnectionPorperties []ValueDomain
	DeploymentZoneRule   string
	//EndpointImplementationDetails DeploymentZoneEndpoint `json:"EndpointImplementationDetails.DeploymentZoneEndpoint"`
	EndpointImplementationDetails struct {
		DeploymentZoneEndpoint
	}
	EndpointKey        string
	ImplementationCode string
	URI                string `json:"Uri"`
}

// DeploymentZoneEndpoint contains information about a Deployment Zone Endpoint
type DeploymentZoneEndpoint struct {
	BindingQName     string
	BindingType      string
	ContainerKey     string
	DeploymentZoneID string
	EndpointHostname string
	EndpointKey      string
	EndpointPath     string
	GatewayHostName  string
	GatewayHostPath  string
	ListenerName     string
	Path             string
	Protocol         string
	Public           bool
	URL              string `json:"Url"`
}

// APICreatedResponse is the information that comes back from a successfully created API
type APICreatedResponse struct {
	APIID           string
	Name            string
	Description     string
	Visibility      string
	LatestVersionID string
	IsFollowed      bool
	RatingSummary   RatingSummary
	APIVersion      APIVersion
	AdminGroupID    string
	Created         string
	Updated         string
	AvatarURL       string
}

// RatingSummary holds a summary of ratings for an API
type RatingSummary struct {
	One   int
	Two   int
	Three int
	Four  int
	Five  int
}

// APIVersion contains information about a version of an API
type APIVersion struct {
	APIVersionID       string        `json:"APIVersionID"`
	APIID              string        `json:"APIID"`
	Name               string        `json:"Name"`
	Description        string        `json:"Description"`
	Tag                []interface{} `json:"Tag"`
	ProductionEndpoint string        `json:"ProductionEndpoint"`
	Endpoints          struct {
		Endpoint []struct {
			CName                string `json:"CName"`
			Category             string `json:"Category"`
			URI                  string `json:"Uri"`
			DeploymentZoneRule   string `json:"DeploymentZoneRule"`
			ConnectionProperties []struct {
				Name  string `json:"Name"`
				Value string `json:"Value"`
			} `json:"ConnectionProperties"`
			BindingQName                  string `json:"BindingQName"`
			BindingType                   string `json:"BindingType"`
			EndpointKey                   string `json:"EndpointKey"`
			EndpointImplementationDetails struct {
				DeploymentZoneEndpoint struct {
					DeploymentZoneID string `json:"DeploymentZoneID"`
					EndpointKey      string `json:"EndpointKey"`
					ListenerName     string `json:"ListenerName"`
					ContainerKey     string `json:"ContainerKey"`
					GatewayHostName  string `json:"GatewayHostName"`
					GatewayHostPath  string `json:"GatewayHostPath"`
					EndpointHostName string `json:"EndpointHostName"`
					EndpointPath     string `json:"EndpointPath"`
					Protocol         string `json:"Protocol"`
					Path             string `json:"Path"`
					URL              string `json:"Url"`
					BindingQName     string `json:"BindingQName"`
					BindingType      string `json:"BindingType"`
					Public           bool   `json:"Public"`
				} `json:"DeploymentZoneEndpoint"`
			} `json:"EndpointImplementationDetails"`
			ImplementationCode string `json:"ImplementationCode"`
		} `json:"Endpoint"`
	} `json:"Endpoints"`
	Visibility                           string    `json:"Visibility"`
	Created                              time.Time `json:"Created"`
	Updated                              time.Time `json:"Updated"`
	State                                string    `json:"State"`
	ProductionEndpointAccessAutoApproved bool      `json:"ProductionEndpointAccessAutoApproved"`
	SandboxEndpointAccessAutoApproved    bool      `json:"SandboxEndpointAccessAutoApproved"`
	RatingSummary                        RatingSummary
	SandboxAnonymousAccessAllowed        bool   `json:"SandboxAnonymousAccessAllowed"`
	ProductionAnonymousAccessAllowed     bool   `json:"ProductionAnonymousAccessAllowed"`
	ResourceLevelPermissionsSupported    bool   `json:"ResourceLevelPermissionsSupported"`
	APIOwnedImplementations              bool   `json:"APIOwnedImplementations"`
	ProductionServiceKey                 string `json:"ProductionServiceKey"`
	APIDesign                            struct {
		CommonDesign bool
	}
}

// APIDetails is the response to an API details request
type APIDetails struct {
	APIID           string `json:"APIID"`
	Name            string `json:"Name"`
	Description     string `json:"Description"`
	Visibility      string `json:"Visibility"`
	LatestVersionID string `json:"LatestVersionID"`
	IsFollowed      bool   `json:"IsFollowed"`
	RatingSummary   RatingSummary
	APIVersion      APIVersion
	AdminGroupID    string    `json:"AdminGroupID"`
	Created         time.Time `json:"Created"`
	Updated         time.Time `json:"Updated"`
	AvatarURL       string    `json:"AvatarURL"`
}
