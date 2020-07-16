package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/f4nt0md3v/tic-tac-go-beeline/app/constants"
)

// Use godot package to load/read the .env file and
// return the value of the key
func GetEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

type Mode int

const (
	Development Mode = iota
	Production
)

func GetAppMode() Mode {
	switch GetEnvVariable(constants.AppMode) {
	case "production":
		return Production
	case "development":
		return Development
	default:
		return Development
	}
}
