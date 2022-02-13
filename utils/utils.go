package utils

import (
	"os"
	"path/filepath"
)

// DeskifyPath returns the path of deskify folder.
func DeskifyPath() string {
	homedir, _ := os.UserHomeDir()
	return filepath.Join(homedir, ".deskify")
}

// SqlitePath returns the path of sqlite file.
func SqlitePath() string {
	homedir, _ := os.UserHomeDir()
	deskifyPath := filepath.Join(homedir, ".deskify")

	return deskifyPath + "/gorm.db"
}

// DotEnvPath returns the path of .env file.
func DotEnvPath() string {
	deskifyPath := DeskifyPath()
	dotEnvPath := deskifyPath + "/.env"

	if _, err := os.Stat(dotEnvPath); os.IsNotExist(err) {
		os.MkdirAll(deskifyPath, os.ModePerm)
		file, _ := os.Create(dotEnvPath)
		file.Close()
	}

	return dotEnvPath
}
