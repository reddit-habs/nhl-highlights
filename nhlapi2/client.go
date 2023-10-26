package nhlapi2

import (
	"context"
	"fmt"

	"github.com/carlmjohnson/requests"
)

const (
	nhlBrightCoveAccount   = "6415718365001"
	nhlBrightCovePolicyKey = "BCpkADawqM3l37Vq8trLJ95vVwxubXYZXYglAopEZXQTHTWX3YdalyF9xmkuknxjBgiMYwt8VZ_OZ1jAjYxz_yzuNh_cjC3uOaMspVTD-hZfNUHtNnBnhVD0Gmsih8TBF8QlQFXiCQM3W_u4ydJ1qK2Rx8ZutCUg3PHb7Q"
)

type Client struct {
}

func NewClient() Client {
	return Client{}
}

func (Client) Schedule(date string) (*ScheduleResponse, error) {
	r := ScheduleResponse{}
	if err := requests.URL(fmt.Sprintf("https://api-web.nhle.com/v1/schedule/%s", date)).ToJSON(&r).Fetch(context.Background()); err != nil {
		return nil, err
	}
	return &r, nil
}

func (Client) Landing(gameID int64) (*LandingResponse, error) {
	r := LandingResponse{}
	url := fmt.Sprintf("https://api-web.nhle.com/v1/gamecenter/%d/landing", gameID)
	err := requests.URL(url).
		ToJSON(&r).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (Client) VideoMetadata(videoID int64) (*VideoMetadataResponse, error) {
	r := VideoMetadataResponse{}
	url := fmt.Sprintf("https://edge.api.brightcove.com/playback/v1/accounts/%s/videos/%d", nhlBrightCoveAccount, videoID)
	err := requests.URL(url).
		Accept(fmt.Sprintf("application/json;pk=%s", nhlBrightCovePolicyKey)).
		ToJSON(&r).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}
	return &r, nil
}
