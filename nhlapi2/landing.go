package nhlapi2

type LandingResponse struct {
	Summary *LandingSummary `json:"summary"`
}

type LandingSummary struct {
	GameVideo *LandingGameVideo `json:"gameVideo"`
}

type LandingGameVideo struct {
	ThreeMinRecap int64 `json:"threeMinRecap"`
	CondensedGame int64 `json:"condensedGame"`
}
