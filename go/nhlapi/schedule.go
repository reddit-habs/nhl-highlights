package nhlapi

type ScheduleResponse struct {
	Dates []*ScheduleDate `json:"dates"`
}

type ScheduleDate struct {
	Date  string          `json:"date"`
	Games []*ScheduleGame `json:"games"`
}

type ScheduleGame struct {
	GameID int64          `json:"gamePk"`
	Type   string         `json:"gameType"`
	Season string         `json:"season"`
	Teams  *ScheduleTeams `json:"teams"`
}

type ScheduleTeams struct {
	Away *ScheduleTeam `json:"away"`
	Home *ScheduleTeam `json:"home"`
}

type ScheduleTeam struct {
	TeamID *ScheduleTeamID `json:"team"`
}

type ScheduleTeamID struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}
