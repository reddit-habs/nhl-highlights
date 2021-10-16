package generate

import (
	"bufio"
	_ "embed"
	"os"
	"path"
	"sort"
	"text/template"

	"github.com/sbstp/nhl-highlights/nhlapi"
	"github.com/sbstp/nhl-highlights/repository"
)

//go:embed template.html
var tplSource string

var tpl = template.Must(template.New("template.html").Parse(tplSource))

func Pages(outputDir string, games []*repository.Game) error {
	bySeason := groupBySeason(games)

	seasons := make([]string, 0)
	for season := range bySeason {
		seasons = append(seasons, season)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(seasons)))

	for season, games := range bySeason {
		outputPath := path.Join(outputDir, season, "index.html")

		err := generateFile(outputPath, &Root{
			Season:  season,
			Seasons: seasons,
			Teams:   nhlapi.Teams,
			Dates:   groupByDate(games),
		})
		if err != nil {
			return err
		}

		teams := groupByTeam(games)
		for team, games := range teams {
			outputPath := path.Join(outputDir, season, team+".html")

			err := generateFile(outputPath, &Root{
				Season:  season,
				Seasons: seasons,
				Teams:   nhlapi.Teams,
				Dates:   groupByDate(games),
			})
			if err != nil {
				return err
			}
		}
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

type Root struct {
	Season  string
	Seasons []string
	Teams   []*nhlapi.Team
	Dates   []Date
}

type Date struct {
	Date  string
	Games []*repository.Game
}

func cleanupSeason(s string) string {
	return s[:4] + "-" + s[4:]
}

func groupBySeason(games []*repository.Game) map[string][]*repository.Game {
	result := make(map[string][]*repository.Game)
	for _, game := range games {
		season := cleanupSeason(game.Season)
		result[season] = append(result[season], game)
	}
	return result
}

func groupByDate(games []*repository.Game) []Date {
	temp := make(map[string][]*repository.Game)
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

func groupByTeam(games []*repository.Game) map[string][]*repository.Game {
	result := make(map[string][]*repository.Game)
	for _, game := range games {
		result[game.Away] = append(result[game.Away], game)
		result[game.Home] = append(result[game.Home], game)
	}
	return result
}
