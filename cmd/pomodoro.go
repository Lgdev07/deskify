package cmd

import (
	"fmt"
	"log"
	"sync"

	"github.com/Lgdev07/deskify/services/pomodoro"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

func InitPomodoroCmd(db *gorm.DB) {
	var (
		focusTimer int
		restTimer  int
	)

	var cmdPomodoro = &cobra.Command{
		Use:   "pomodoro [action]",
		Short: "Configure pomodoro technique and be notified",
		Long:  "Configure pomodoro technique and be notified",
		Run: func(cmd *cobra.Command, args []string) {
			if focusTimer == 0 || restTimer == 0 {
				fmt.Println("Focus or rest cannot be 0")
				return
			}
			SaveInformation(db, focusTimer, restTimer)
		},
	}

	var cmdPomodoroRun = &cobra.Command{
		Use:   "run",
		Short: "Start pomodoro schedule",
		Long:  "Command to start pomodoro schedule base on your configuration.",
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup

			pomodoro.Initialize(&wg, db)

			wg.Wait()
		},
	}

	cmdPomodoro.Flags().IntVarP(&focusTimer, "focus", "f", 0, "focus timer")
	cmdPomodoro.Flags().IntVarP(&restTimer, "rest", "r", 0, "rest timer")
	cmdPomodoro.MarkFlagRequired("focus")
	cmdPomodoro.MarkFlagRequired("rest")

	rootCmd.AddCommand(cmdPomodoro)

	cmdPomodoro.AddCommand(cmdPomodoroRun)

}

func SaveInformation(db *gorm.DB, focusTimer, restTimer int) {
	pomodoroValue := &pomodoro.Pomodoro{}
	db.Model(&pomodoro.Pomodoro{}).Find(pomodoroValue)

	if pomodoroValue.Focus == 0 {
		CreatePomodoro(db, focusTimer, restTimer)
		return
	}

	pomodoroValue.Focus = focusTimer
	pomodoroValue.Rest = restTimer
	db.Save(&pomodoroValue)

}

func CreatePomodoro(db *gorm.DB, focusTimer, restTimer int) {
	newPomodoro := &pomodoro.Pomodoro{
		Focus: focusTimer,
		Rest:  restTimer,
	}

	err := db.Model(&pomodoro.Pomodoro{}).Create(newPomodoro).Error
	if err != nil {
		log.Fatal(err)
	}

}
