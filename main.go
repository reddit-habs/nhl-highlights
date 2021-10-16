package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sbstp/nhl-highlights/addrof"
	"github.com/sbstp/nhl-highlights/generate"
	"github.com/sbstp/nhl-highlights/nhlapi"
	"github.com/sbstp/nhl-highlights/repository"
)

var incremental bool
var startDate string
var endDate string
var outputDir string

func main() {
	flag.BoolVar(&incremental, "incremental", false, "try to get highlights from the past few days")
	flag.StringVar(&startDate, "start-date", "", "start date of scan")
	flag.StringVar(&endDate, "end-date", "", "end date of scan")
	flag.StringVar(&outputDir, "output-dir", "output", "directory where HTML files end up")
	flag.Parse()

	if !incremental && (len(startDate) == 0 || len(endDate) == 0) {
		fmt.Fprintln(os.Stderr, "-start-date and -end-date are required in scan mode")
		return
	}

	if incremental && (len(startDate) > 0 || len(endDate) > 0) {
		fmt.Fprintln(os.Stderr, "-start-date and -end-date have no meaning in incremental mode")
	}

	if err := realMain(); err != nil {
		log.Fatal(err)
	}
}

func realMain() error {
	repo, err := repository.New("games.db")
	if err != nil {
		return err
	}
	client := nhlapi.NewClient()

	schedule, err := client.Schedule(startDate, endDate)
	if err != nil {
		return err
	}

	for _, date := range schedule.Dates {
		for _, game := range date.Games {
			if game.Type == "A" || game.Type == "WA" {
				continue
			}
			exists, err := repo.GetGame(game.GameID)
			if err != nil {
				return err
			}
			if exists == nil {
				g := gameFromSchedule(date.Date, game)
				if g == nil {
					continue
				}
				if err := repo.UpsertGame(g); err != nil {
					return err
				}
			}
		}
	}

	missing, err := repo.GetGamesMissingContent(incremental)
	for _, game := range missing {
		log.Printf("Getting content for game %d, date=%s", game.GameID, game.Date)
		content, err := client.Content(game.GameID)
		if err != nil {
			return err
		}
		for _, epg := range content.Media.EPG {
			if epg.Title == "Recap" && len(epg.Items) > 0 {
				game.Recap = addrof.String(getBestMp4Playback(epg.Items[0].Playbacks))
			}
			if epg.Title == "Extended Highlights" && len(epg.Items) > 0 {
				game.Extended = addrof.String(getBestMp4Playback(epg.Items[0].Playbacks))
			}
		}
		if err := repo.UpsertGame(game); err != nil {
			return err
		}
	}

	games, err := repo.GetGames()
	if err != nil {
		return err
	}

	if err := generate.Pages(outputDir, games); err != nil {
		return err
	}

	return nil
}

func gameFromSchedule(date string, game *nhlapi.ScheduleGame) *repository.Game {
	away, ok := nhlapi.TeamsByID[game.Teams.Away.TeamID.ID]
	if !ok {
		return nil
	}
	home, ok := nhlapi.TeamsByID[game.Teams.Home.TeamID.ID]
	if !ok {
		return nil
	}
	return &repository.Game{
		GameID:   game.GameID,
		Date:     date,
		Type:     game.Type,
		Away:     away.Abbrev,
		Home:     home.Abbrev,
		Season:   game.Season,
		Recap:    nil,
		Extended: nil,
	}
}

func getBestMp4Playback(playbacks []*nhlapi.ContentPlayback) string {
	links := make([]string, 0, 4)
	for _, pb := range playbacks {
		if strings.HasSuffix(pb.URL, ".mp4") {
			links = append(links, pb.URL)
		}
	}
	return links[len(links)-1]
}
