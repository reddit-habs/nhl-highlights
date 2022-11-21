package generate

import (
	"bufio"
	_ "embed"
	"html/template"
	"log"
	"os"
	"path"
	"sort"
	"time"

	"github.com/sbstp/nhl-highlights/models"
	"github.com/sbstp/nhl-highlights/nhlapi"
)

//go:embed template.html
var tplSource string

var tpl = template.Must(template.New("template.html").Parse(tplSource))

func Pages(teamsCache *nhlapi.TeamsCache, outputDir string, games []*models.Game) error {
	bySeason := groupBySeason(games)

	seasons := make([]string, 0)
	for season := range bySeason {
		seasons = append(seasons, season)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(seasons)))

	for season, games := range bySeason {
		outputPath := path.Join(outputDir, season, "index.html")
		teams := teamsForSeason(teamsCache, games)

		err := generateFile(outputPath, NewRoot(season, seasons, games, teams))
		if err != nil {
			return err
		}

		for team, games := range groupByTeam(games) {
			outputPath := path.Join(outputDir, season, team+".html")

			err := generateFile(outputPath, NewRoot(season, seasons, games, teams))
			if err != nil {
				return err
			}
		}
	}

	// Remove "current" symlink if it already exists.
	// os.Symlink won't overwrite it.
	currentPath := path.Join(outputDir, "current")
	if err := os.Remove(currentPath); err != nil {
		log.Printf("Could not remove current: %v", err)
	}

	if err := os.Symlink(path.Join(seasons[0]), currentPath); err != nil {
		log.Printf("Error symlinking: %v", err)
	}

	return nil
}

func generateFile(outputPath string, root *Root) error {
	os.MkdirAll(path.Dir(outputPath), 0755)

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	if err := tpl.Execute(writer, root); err != nil {
		return err
	}

	return nil
}

func NewRoot(season string, seasons []string, games []*models.Game, teams []*nhlapi.Team) *Root {
	return &Root{
		Season:         season,
		Seasons:        seasons,
		Teams:          teams,
		Dates:          groupByDate(games),
		GenerationDate: time.Now().Local().Format("2006-01-02 03:04:05"),
	}
}

type Root struct {
	Season         string
	Seasons        []string
	Teams          []*nhlapi.Team
	Dates          []Date
	GenerationDate string
}

type Date struct {
	Date  string
	Games []*models.Game
}

func cleanupSeason(s string) string {
	return s[:4] + "-" + s[4:]
}

func groupBySeason(games []*models.Game) map[string][]*models.Game {
	result := make(map[string][]*models.Game)
	for _, game := range games {
		season := cleanupSeason(game.Season)
		result[season] = append(result[season], game)
	}
	return result
}

func groupByDate(games []*models.Game) []Date {
	temp := make(map[string][]*models.Game)
	for _, game := range games {
		temp[game.Date] = append(temp[game.Date], game)
	}
	result := make([]Date, 0, len(temp))
	for date, games := range temp {
		result = append(result, Date{
			Date:  date,
			Games: games,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date > result[j].Date
	})
	return result
}

func groupByTeam(games []*models.Game) map[string][]*models.Game {
	result := make(map[string][]*models.Game)
	for _, game := range games {
		result[game.Away] = append(result[game.Away], game)
		result[game.Home] = append(result[game.Home], game)
	}
	return result
}

func teamsForSeason(teamsCache *nhlapi.TeamsCache, games []*models.Game) []*nhlapi.Team {
	temp := make(map[string]struct{})
	for _, game := range games {
		temp[game.Away] = struct{}{}
		temp[game.Home] = struct{}{}
	}
	teams := make([]*nhlapi.Team, 0, len(temp))
	for team := range temp {
		teams = append(teams, teamsCache.GetByAbbrev(team))
	}
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Abbrev < teams[j].Abbrev
	})
	return teams
}
