package nhlapi2

const (
	GameTypePreseason = 1
	GameTypeRegular   = 2
	GameTypePlayoffs  = 3
)

type ScheduleResponse struct {
	NextStartDate     string          `json:"nextStartDate"`
	PreviousStartDate string          `json:"previousStartDate"`
	GameWeek          []*ScheduleDate `json:"gameWeek"`
}

type ScheduleDate struct {
	Date  string          `json:"date"`
	Games []*ScheduleGame `json:"games"`
}

type ScheduleGame struct {
	GameID        int64         `json:"id"`
	Season        int64         `json:"season"`
	AwayTeam      *ScheduleTeam `json:"awayTeam"`
	HomeTeam      *ScheduleTeam `json:"homeTeam"`
	ThreeMinRecap string        `json:"threeMinRecap"`
	Type          int32         `json:"gameType"`
}

type ScheduleTeam struct {
	ID     int32  `json:"id"`
	Abbrev string `json:"abbrev"`
}
