package nhlapi

type ContentResponse struct {
	Media ContentMedia `json:"media"`
}

type ContentMedia struct {
	EPG []*ContentEPG `json:"epg"`
}

type ContentEPG struct {
	Title string         `json:"title"`
	Items []*ContentItem `json:"items"`
}

type ContentItem struct {
	Type      string             `json:"type"`
	Playbacks []*ContentPlayback `json:"playbacks"`
}

type ContentPlayback struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
