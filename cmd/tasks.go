package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Lgdev07/deskify/services/tasks"
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
		Use:   "rem [task]",
		Short: "Remove a task",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			task := fmt.Sprintf(strings.Join(args, " "))
			TaskRem(db, task)
		},
	}

	var cmdTaskDone = &cobra.Command{
		Use:   "done [task]",
		Short: "Mark a task as done",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			task := fmt.Sprintf(strings.Join(args, " "))
			TaskDone(db, task)
		},
	}

	var cmdTaskList = &cobra.Command{
		Use:   "list",
		Short: "Show all active tasks",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			TaskListActive(db)
		},
	}

	cmdTaskAdd.Flags().IntVarP(&timer, "timer", "t", 0, "timer to be remebered every x minutes")
	cmdTaskAdd.MarkFlagRequired("timer")

	rootCmd.AddCommand(cmdTask)

	cmdTask.AddCommand(cmdTaskAdd)
	cmdTask.AddCommand(cmdTaskRem)
	cmdTask.AddCommand(cmdTaskDone)
	cmdTask.AddCommand(cmdTaskList)

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

func TaskRem(db *gorm.DB, taskName string) {
	task := tasks.Task{}

	db.Model(&tasks.Task{}).Where("name = ?", taskName).First(&task)

	if task.Name == "" {
		fmt.Println("We did not find a task with that name")
		return
	}

	err := db.Model(&tasks.Task{}).Where("name = ?", taskName).Delete(&tasks.Task{}).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Task %s Deleted with success\n", taskName)

}

func TaskDone(db *gorm.DB, taskName string) {
	task := tasks.Task{}

	db.Model(&tasks.Task{}).Where("name = ?", taskName).First(&task)

	if task.Name == "" {
		fmt.Println("We did not find a task with that name")
		return
	}

	err := db.Model(&tasks.Task{}).Where("name = ?", taskName).Update("is_done", true).Error
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task %s Marked as Done!\n", taskName)
}

func TaskListActive(db *gorm.DB) {
	taskList := []tasks.Task{}

	db.Model(&tasks.Task{}).Where("is_done = 0").Find(&taskList)

	if len(taskList) == 0 {
		fmt.Println("No active tasks found")
		return
	}

	for _, value := range taskList {
		fmt.Printf("Name: %s, Timer: %d\n", value.Name, value.BeRememberedMinutes)
	}
}
