package cm

type ApisResponse struct {
	Channel Channel `json:"channel"`
}

type Channel struct {
	Title string `json:"title"`
	Items []Item `json:"item"`
}

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

type ValueDomain struct {
	Value  string `json:"value"`
	Domain string `json:"domain"`
}

type Guid struct {
	Value string `json:"value"`
}

type EntityReference struct {
	Title    string
	Guid     string
	Category ValueDomain
}
