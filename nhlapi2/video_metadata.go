package nhlapi2

type VideoMetadataResponse struct {
	Sources []*VideoMetadataSource `json:"sources"`
}

type VideoMetadataSource struct {
	AvgBitrate int64  `json:"avg_bitrate"`
	Codec      string `json:"codec"`
	Container  string `json:"container"`
	Height     int32  `json:"height"`
	Width      int32  `json:"width"`
	Src        string `json:"src"`
}
