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
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	Category         []ValueDomain     `json:"category"`
	Guid             Guid              `json:"guid"`
	PubDate          string            `json:"pubDate"`
	EntityReferences []EntityReference `json:"EntityReferences.EntityReference"`
	EntityReference  EntityReference   `json:"EntityReference"`
	ImageUrl         string            `json:"Image.Url"`
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
	Guid     string
	Category ValueDomain
}
