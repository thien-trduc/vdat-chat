package utils

import (
	"github.com/joho/godotenv"
	_ "log"
	os "os"
)

func LoadEnv(fileName string) bool {
	if len(fileName) == 0 {
		fileName = "dev.env"
	}

	err := godotenv.Load(fileName)

	return err != nil
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
