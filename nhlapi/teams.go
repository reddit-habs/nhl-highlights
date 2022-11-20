package nhlapi

import (
	"sort"
)

type TeamsResponse struct {
	Teams []*Team `json:"teams"`
}

type Team struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Abbrev   string `json:"abbreviation"`
	TeamName string `json:"teamName"`
	Location string `json:"locationName"`
}

type TeamsCache struct {
	Teams         []*Team
	teamsByID     map[int32]*Team
	teamsByAbbrev map[string]*Team
}

func NewTeamsCache(c Client) (*TeamsCache, error) {
	response, err := c.Teams()
	if err != nil {
		return nil, err
	}

	teams := make([]*Team, 0, len(response.Teams))
	teamsByID := make(map[int32]*Team)
	teamsByAbbrev := make(map[string]*Team)

	for _, team := range response.Teams {
		teams = append(teams, team)
		teamsByID[team.ID] = team
		teamsByAbbrev[team.Abbrev] = team
	}

	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Abbrev < teams[j].Abbrev
	})

	return &TeamsCache{
		Teams:         teams,
		teamsByID:     teamsByID,
		teamsByAbbrev: teamsByAbbrev,
	}, nil
}

func (t *TeamsCache) GetByID(id int32) (*Team, bool) {
	v, ok := t.teamsByID[id]
	return v, ok
}

func (t *TeamsCache) GetByAbbrev(abbrev string) *Team {
	return t.teamsByAbbrev[abbrev]
}
