package cm

// ApisResponse is the main struct for the RSS feed
type ApisResponse struct {
	Channel Channel `json:"channel"`
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
