package database

import (
	"log"
	"os"

	"github.com/Lgdev07/deskify/services/pomodoro"
	"github.com/Lgdev07/deskify/services/tasks"
	"github.com/Lgdev07/deskify/services/twitch"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

type Database struct {
	DB *gorm.DB
}

func (s *Database) Initialize() {
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error getting env, %v", err)
		}
	}

	s.DB, _ = gorm.Open("sqlite3", "./gorm.db")

	s.DB.AutoMigrate(
		&twitch.Twitch{},
		&tasks.Task{},
		&pomodoro.Pomodoro{},
	)

}
