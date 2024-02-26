package helpers

import (
	"fmt"
	"os"
	"strconv"
)

func GetIntValueFromEnv(key string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		fmt.Println("Fail to convert "+key+", set default value to:"+strconv.Itoa(defaultValue), err)
		value = defaultValue
	}
	return value
}
