package nhlapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	c *http.Client
}

func NewClient() Client {
	return Client{
		c: http.DefaultClient,
	}
}

func (c Client) Schedule(startDate, endDate string) (*ScheduleResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://statsapi.web.nhl.com/api/v1/schedule?startDate=%s&endDate=%s", startDate, endDate), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	body := &ScheduleResponse{}
	if err := json.NewDecoder(resp.Body).Decode(body); err != nil {
		return nil, err
	}

	return body, err
}
