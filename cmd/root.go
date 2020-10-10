package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/Lgdev07/deskify/database"
	"github.com/Lgdev07/deskify/tasks"
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
			var wg sync.WaitGroup

			twitch.Initialize(&wg, db.DB)
			tasks.Initialize(&wg, db.DB)

			wg.Wait()
		},
	}

	rootCmd.AddCommand(cmdRun)
	InitTwitchCmd(db.DB)
	InitTasksCmd(db.DB)

}
