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

func (c Client) doGet(url string, response any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("error performing GET %q: got status %d", url, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return err
	}

	return nil
}

func (c Client) Schedule(startDate, endDate string) (*ScheduleResponse, error) {
	response := &ScheduleResponse{}

	if err := c.doGet(fmt.Sprintf("https://statsapi.web.nhl.com/api/v1/schedule?startDate=%s&endDate=%s", startDate, endDate), response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c Client) Content(gameID int64) (*ContentResponse, error) {
	response := &ContentResponse{}

	if err := c.doGet(fmt.Sprintf("https://statsapi.web.nhl.com/api/v1/game/%d/content", gameID), response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c Client) Teams() (*TeamsResponse, error) {
	response := &TeamsResponse{}

	if err := c.doGet("https://statsapi.web.nhl.com/api/v1/teams", response); err != nil {
		return nil, err
	}

	return response, nil
}
