package pomodoro

import (
	"fmt"
	"sync"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/jinzhu/gorm"
)

type Pomodoro struct {
	gorm.Model
	Focus int `gorm:"default:0" json:"focus"`
	Rest  int `gorm:"default:0" json:"rest"`
}

func Initialize(wg *sync.WaitGroup, db *gorm.DB) {
	pomodoro := &Pomodoro{}

	db.Model(&Pomodoro{}).Find(pomodoro)

	if pomodoro.Focus == 0 || pomodoro.Rest == 0 {
		return
	}

	wg.Add(1)

	go func() {
		for {
			message := fmt.Sprintf("Focus time started, it will remain for %d minutes", pomodoro.Focus)
			err := beeep.Notify(message, "", "assets/pomodoroFocus.png")
			if err != nil {
				panic(err)
			}

			time.Sleep(time.Duration(pomodoro.Focus) * time.Minute)

			message = fmt.Sprintf("Focus time finished, now you will rest for %d minutes", pomodoro.Rest)
			err = beeep.Notify(message, "", "assets/pomodoroRest.png")
			if err != nil {
				panic(err)
			}

			time.Sleep(time.Duration(pomodoro.Rest) * time.Minute)

		}
	}()

}
