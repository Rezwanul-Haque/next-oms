package cmd

import (
	server "next-oms/app/http"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: serve,
}

func serve(cmd *cobra.Command, args []string) {
	// http server start
	server.Start()
}
