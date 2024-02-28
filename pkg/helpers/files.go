package helpers

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func GetIntValueFromEnv(key string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		fmt.Println("Fail to convert "+key+", set default value to:"+strconv.Itoa(defaultValue), err)
		value = defaultValue
	}
	return value
}

func LoadEnvFile() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
}
