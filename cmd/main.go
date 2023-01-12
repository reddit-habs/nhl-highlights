package main

import (
	"log"

	"github.com/sbstp/nhl-highlights/repository"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "nhl-highlights",
		Short: "Archive & display NHL highlights",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	rootCmd.AddCommand(newArchiveIncrementalCmd())
	rootCmd.AddCommand(newArchiveRangeCmd())
	rootCmd.AddCommand(newServeCmd())
	rootCmd.AddCommand(newCreateDatabaseCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func newArchiveIncrementalCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "archive-incremental",
		Short: "archive higlights from the past few days",
		RunE: func(cmd *cobra.Command, args []string) error {
			return realMain(true, "", "")
		},
	}
}

func newArchiveRangeCmd() *cobra.Command {
	var startDate string
	var endDate string

	cmd := &cobra.Command{
		Use:   "archive-range",
		Short: "archive higlights from the given date range",
		RunE: func(cmd *cobra.Command, args []string) error {
			return realMain(false, startDate, endDate)
		},
	}

	cmd.Flags().StringVar(&startDate, "start-date", "", "start date of scan")
	cmd.Flags().StringVar(&endDate, "end-date", "", "end date of scan")
	cmd.MarkFlagRequired("start-date")
	cmd.MarkFlagRequired("end-date")

	return cmd
}

func newServeCmd() *cobra.Command {
	var bindAddress string

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "serve highlights as HTML",
		RunE: func(cmd *cobra.Command, args []string) error {
			return serve(bindAddress)
		},
	}

	cmd.Flags().StringVar(&bindAddress, "bind-address", ":9999", "bind address and port to use")

	return cmd
}

func newCreateDatabaseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-database",
		Short: "Initialize an empty database at the given path",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dbpath := args[0]
			repo, err := repository.New(dbpath)
			if err != nil {
				return err
			}
			if err = repo.Close(); err != nil {
				return err
			}
			return nil
		},
	}
}

// func scanHighlights(repo *repository.Repository, content *nhlapi.ContentResponse, game *models.Game) {
// 	var highlights []*models.Highlight

// 	for _, item := range content.Highlights.Scoreboard.Items {
// 		video := getBestMp4Playback(item.Playbacks)

// 		var eventID null.Int64
// 		for _, kw := range item.Keywords {
// 			if kw.Type == "statsEventId" {
// 				eventID = null.Int64From(stringToInt64(kw.Value))
// 			}
// 		}

// 		highlights = append(highlights, &models.Highlight{
// 			ID:          stringToInt64(item.ID),
// 			GameID:      game.GameID,
// 			MediaURL:    video,
// 			EventID:     eventID,
// 			Title:       null.StringFrom(item.Title),
// 			Blurb:       null.StringFrom(item.Blurb),
// 			Description: null.StringFrom(item.Description),
// 		})
// 	}

// 	repo.UpsertHighlights(highlights)

// }

// func stringToInt64(s string) int64 {
// 	x, err := strconv.ParseInt(s, 10, 64)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return x
// }
