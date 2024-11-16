package cmd

import (
	"fmt"
	"next-oms/infra/config"
	"next-oms/infra/conn/cache"
	"next-oms/infra/conn/db"
	"next-oms/infra/logger"
	"os"

	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "next-oms",
		Short: "implementing oms for next",
	}
)

func init() {
	RootCmd.AddCommand(serveCmd)
}

// Execute executes the root command
func Execute() {
	config.LoadConfig()
	logger.NewLogClient(config.App().LogLevel)
	lc := logger.Client()
	db.NewDbClient(lc)
	cache.NewCacheClient(lc)

	lc.Info("about to start the application")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
