package cmd

import (
	"context"
	"go-delayed-job/cmd/playground"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go-delayed-job/cmd/job"

	"github.com/spf13/cobra"
)

func Start() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	ctx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		cancel()
	}()

	rootCmd := &cobra.Command{}
	cmd := []*cobra.Command{
		{
			Use:   "playground",
			Short: "Run Playground",
			Run: func(cmd *cobra.Command, args []string) {
				playground.RunPlayground(ctx)
			},
		},
		{
			Use:   "job:expire-order",
			Short: "Run Job Expire Order",
			Run: func(cmd *cobra.Command, args []string) {
				job.RunJobExpireOrder(ctx)
			},
		},
	}

	rootCmd.AddCommand(cmd...)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
