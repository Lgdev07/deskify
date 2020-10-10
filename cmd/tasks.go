package cmd

import (
	"fmt"
	"strings"

	"github.com/Lgdev07/deskify/tasks"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

func InitTasksCmd(db *gorm.DB) {
	var timer int

	var cmdTask = &cobra.Command{
		Use:   "task [action]",
		Short: "Do an action with the command task",
		Long:  "task command preceed an action.",
		Args:  cobra.MinimumNArgs(1),
	}

	var cmdTaskAdd = &cobra.Command{
		Use:   "add [task] --timer [minutes]",
		Short: "Add a task to be notified every x minutes",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			task := fmt.Sprintf(strings.Join(args, " "))
			TaskAdd(db, task, timer)

		},
	}

	var cmdTaskRem = &cobra.Command{
		Use:   "add [task] --timer [minutes]",
		Short: "Add a task to be notified every x minutes",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for i := 0; i < timer; i++ {
				fmt.Println("Echo: " + strings.Join(args, " "))
			}
		},
	}

	cmdTaskAdd.Flags().IntVarP(&timer, "timer", "t", 0, "timer to be remebered every x minutes")
	cmdTaskAdd.MarkFlagRequired("timer")

	rootCmd.AddCommand(cmdTask)

	cmdTask.AddCommand(cmdTaskAdd)
	cmdTask.AddCommand(cmdTaskRem)

}

func TaskAdd(db *gorm.DB, taskName string, timer int) {
	task := tasks.Task{}

	db.Model(&tasks.Task{}).Where("name = ?", taskName).First(&task)

	if task.Name != "" {
		fmt.Println("There is already a task with the same name")
		return
	}

	newTask := &tasks.Task{
		Name:                taskName,
		BeRememberedMinutes: timer,
	}

	db.Create(newTask)
	fmt.Printf("Task %s created with success, you will be remebered every %d minutes\n", taskName, timer)
}
