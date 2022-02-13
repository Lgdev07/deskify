package database

import (
	"log"

	"github.com/Lgdev07/deskify/services/pomodoro"
	"github.com/Lgdev07/deskify/services/tasks"
	"github.com/Lgdev07/deskify/services/twitch"
	"github.com/Lgdev07/deskify/utils"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func (s *Database) Initialize() {
	dotEnvPath := utils.DotEnvPath()

	err := godotenv.Load(dotEnvPath)
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	}

	sqlitePath := utils.SqlitePath()
	s.DB, _ = gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})

	s.DB.AutoMigrate(
		&twitch.Twitch{},
		&tasks.Task{},
		&pomodoro.Pomodoro{},
	)

}
