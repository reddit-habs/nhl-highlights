package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sbstp/nhl-highlights/generate"
	"github.com/sbstp/nhl-highlights/models"
	"github.com/sbstp/nhl-highlights/nhlapi"
	"github.com/sbstp/nhl-highlights/nhlapi2"
	"github.com/sbstp/nhl-highlights/repository"
	"github.com/volatiletech/null/v8"
)

func archive(incremental bool, startDate string, endDate string) error {
	repo, err := repository.New("games.db")
	if err != nil {
		return err
	}
	defer repo.Close()

	client := nhlapi.NewClient()
	teamsCache, err := nhlapi.NewTeamsCache(client)
	if err != nil {
		return err
	}

	clientv2 := nhlapi2.NewClient()

	archiveChunk := func(schedule *nhlapi2.ScheduleResponse) error {
		for _, date := range schedule.GameWeek {
			log.Printf("Date: %s", date.Date)
			for _, game := range date.Games {
				log.Printf("Game: %s at %s", game.AwayTeam.Abbrev, game.HomeTeam.Abbrev)
				if !isGameRelavant(game) {
					continue
				}
				exists, err := repo.GetGame(game.GameID, date.Date)
				if err != nil {
					return err
				}
				if exists == nil {
					log.Printf("Adding game %d on date %s", game.GameID, date.Date)
					g := gameFromSchedule(teamsCache, date.Date, game)
					if g == nil {
						continue
					}
					if err := repo.UpsertGame(g); err != nil {
						return err
					}
				}
			}
		}

		return nil
	}

	if incremental {
		schedule, err := clientv2.Schedule("now")
		if err != nil {
			return err
		}

		if err := archiveChunk(schedule); err != nil {
			return err
		}
	} else {
		cursorDate := startDate
		for {
			log.Printf("Cursor date: %s, end date %s", cursorDate, endDate)
			if cursorDate > endDate {
				break
			}

			schedule, err := clientv2.Schedule(cursorDate)
			if err != nil {
				return err
			}

			if err := archiveChunk(schedule); err != nil {
				return err
			}

			cursorDate = schedule.NextStartDate
		}
	}

	missing, err := repo.GetGamesMissingContent(incremental)
	if err != nil {
		return err
	}

	for _, game := range missing {
		log.Printf("Getting content for game %d, date=%s", game.GameID, game.Date)
		landing, err := clientv2.Landing(game.GameID)
		if err != nil {
			log.Printf("[error] could not get landing: %v", err)
			continue
		}

		if landing.Summary != nil && landing.Summary.GameVideo != nil {
			if landing.Summary.GameVideo.ThreeMinRecap != 0 {
				videoMetadata, err := clientv2.VideoMetadata(landing.Summary.GameVideo.ThreeMinRecap)
				if err != nil {
					log.Printf("[error] could not get video metadata %v", err)
				} else {
					for _, src := range videoMetadata.Sources {
						if src.Codec == "H264" && src.Container == "MP4" && len(src.Src) > 0 {
							game.Recap = null.StringFrom(src.Src)
						}
					}
				}
			}
			if landing.Summary.GameVideo.CondensedGame != 0 {
				videoMetadata, err := clientv2.VideoMetadata(landing.Summary.GameVideo.CondensedGame)
				if err != nil {
					log.Printf("[error] could not get video metadata %v", err)
				} else {
					for _, src := range videoMetadata.Sources {
						if src.Codec == "H264" && src.Container == "MP4" && len(src.Src) > 0 {
							game.Extended = null.StringFrom(src.Src)
						}
					}
				}
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

	stream := make(chan *models.CachedPage)
	ctx, cancel := context.WithCancelCause(context.Background())

	go func() {
		err = generate.Highlights(teamsCache, games, stream)
		defer close(stream)
		if err != nil {
			cancel(err)
			log.Printf("Error generating highlights: %v", err)
		}
	}()

	if err := repo.UpdateCachedPagesIteratively(ctx, stream); err != nil {
		log.Printf("Update iterator error: %v", err)
	}

	log.Print("Archival done.")

	return nil
}

// isGameRelevant checks if the game is relevant for what this program is doing.
// We only care about preseason, regular and playoffs games.
func isGameRelavant(game *nhlapi2.ScheduleGame) bool {
	switch game.Type {
	case nhlapi2.GameTypePreseason, nhlapi2.GameTypeRegular, nhlapi2.GameTypePlayoffs:
		return true
	default:
		return false
	}
}

func gameFromSchedule(teamsCache *nhlapi.TeamsCache, date string, game *nhlapi2.ScheduleGame) *models.Game {
	away := teamsCache.GetByAbbrev(game.AwayTeam.Abbrev)
	home := teamsCache.GetByAbbrev(game.HomeTeam.Abbrev)
	return &models.Game{
		GameID:   game.GameID,
		Date:     date,
		Type:     convertGameType(game.Type),
		Away:     away.Abbrev,
		Home:     home.Abbrev,
		Season:   fmt.Sprintf("%d", game.Season),
		Recap:    null.String{},
		Extended: null.String{},
	}
}

func convertGameType(gameType int32) string {
	switch gameType {
	case nhlapi2.GameTypePreseason:
		return nhlapi.GameTypePreseason
	case nhlapi2.GameTypeRegular:
		return nhlapi.GameTypeRegular
	case nhlapi2.GameTypePlayoffs:
		return nhlapi.GameTypePlayoffs
	}
	panic("unknown game type")
}
