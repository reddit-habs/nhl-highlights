package main

import (
	"context"
	"log"
	"os"
	"os/signal"

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
			return archive(true, "", "")
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
			return archive(false, startDate, endDate)
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
	var incremental bool

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "serve highlights as HTML",
		RunE: func(cmd *cobra.Command, args []string) error {
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, os.Interrupt, os.Interrupt)
			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				<-sigs
				cancel()
			}()
			return serve(ctx, bindAddress, incremental)
		},
	}

	cmd.Flags().StringVar(&bindAddress, "bind-address", ":9999", "bind address and port to use")
	cmd.Flags().BoolVar(&incremental, "incremental", false, "run the incremental archiver periodically")

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
