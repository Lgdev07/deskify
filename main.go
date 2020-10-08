package main

import (
	"github.com/Lgdev07/deskify/cmd"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	cmd.Execute()
}
