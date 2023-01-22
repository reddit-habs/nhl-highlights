package nhlapi

type ContentResponse struct {
	Media      ContentMedia      `json:"media"`
	Highlights ContentHighlights `json:"highlights"`
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

type ContentHighlights struct {
	Scoreboard ContentScoreboard `json:"scoreboard"`
}

type ContentScoreboard struct {
	Items []*ContentHighlightItems `json:"items"`
}

type ContentHighlightItems struct {
	ID          string             `json:"id"`
	Type        string             `json:"type"`
	Title       string             `json:"title"`
	Blurb       string             `json:"blurb"`
	Description string             `json:"description"`
	Keywords    []*ContentKeyword  `json:"keywords"`
	Playbacks   []*ContentPlayback `json:"playbacks"`
}

type ContentKeyword struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
