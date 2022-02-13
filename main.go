package main

import (
	"github.com/Lgdev07/deskify/cmd"

	_ "gorm.io/driver/sqlite"
)

func main() {
	cmd.Execute()
}
