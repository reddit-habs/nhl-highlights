package repository

type Game struct {
	GameID   int64
	Date     string
	Type     string
	Away     string
	Home     string
	Season   string
	Recap    *string
	Extended *string
}
