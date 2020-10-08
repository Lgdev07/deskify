package cmd

import (
	"fmt"
	"os"

	"github.com/Lgdev07/deskify/database"
	"github.com/Lgdev07/deskify/twitch"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	db := database.Database{}
	db.Initialize()

	var cmdRun = &cobra.Command{
		Use:   "run",
		Short: "Initialize the app",
		Long:  "Initialize the app",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			twitch.Initialize(db.DB)
		},
	}

	rootCmd.AddCommand(cmdRun)
	InitTwitchCmd(db.DB)

}
