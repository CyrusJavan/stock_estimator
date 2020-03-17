package cmd

import (
	"github.com/CyrusJavan/stock_estimator/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start an API server",
	Long:  `Start an API server`,
	Run:   runServer,
}

func runServer(cmd *cobra.Command, args []string) {
	s := server.NewAPIServer()

	s.StartServer()
}
