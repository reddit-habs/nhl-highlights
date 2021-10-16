package main

import (
	"log"
	"strings"

	"github.com/sbstp/nhl-highlights/addrof"
	"github.com/sbstp/nhl-highlights/nhlapi"
	"github.com/sbstp/nhl-highlights/repository"
)

func main() {
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

	schedule, err := client.Schedule("2021-10-01", "2021-10-15")
	if err != nil {
		return err
	}

	for _, date := range schedule.Dates {
		for _, game := range date.Games {
			exists, err := repo.GetGame(game.GameID)
			if err != nil {
				return err
			}
			if exists == nil {
				if err := repo.UpsertGame(gameFromSchedule(date.Date, game)); err != nil {
					return err
				}
			}
		}
	}

	missing, err := repo.GetGamesMissingContent()
	for _, game := range missing {
		log.Printf("Getting content for game %d", game.GameID)
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

	return nil
}

func gameFromSchedule(date string, game *nhlapi.ScheduleGame) *repository.Game {
	return &repository.Game{
		GameID:   game.GameID,
		Date:     date,
		Type:     game.Type,
		Away:     nhlapi.TeamsByID[game.Teams.Away.TeamID.ID].Abbrev,
		Home:     nhlapi.TeamsByID[game.Teams.Home.TeamID.ID].Abbrev,
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
