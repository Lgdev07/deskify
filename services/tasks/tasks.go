package tasks

import (
	"log"
	"sync"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/jinzhu/gorm"
)

type Task struct {
	gorm.Model
	Name                string `gorm:"size:100;not null" json:"name"`
	BeRememberedMinutes int    `gorm:"default:0" json:"be_remembered_minutes"`
	IsDone              bool   `gorm:"default:false" json:"is_done"`
}

func Initialize(wg *sync.WaitGroup, db *gorm.DB) {
	tasks, err := getValuesFromDatabase(db)
	if err != nil {
		log.Fatal(err)
	}

	for _, task := range *tasks {
		wg.Add(1)
		go RunFunction(task.Name, task.BeRememberedMinutes)
	}

}

func getValuesFromDatabase(db *gorm.DB) (*[]Task, error) {
	tasks := &[]Task{}

	err := db.Model(&Task{}).Where("is_done = 0").Find(&tasks).Error
	if err != nil {
		return &[]Task{}, err
	}

	return tasks, nil
}

func RunFunction(name string, timer int) {
	for {
		time.Sleep(time.Duration(timer) * time.Minute)
		err := beeep.Notify(name, "Remember to do Your Task!", "assets/task.png")
		if err != nil {
			panic(err)
		}
	}
}
