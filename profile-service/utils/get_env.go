package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func logger(failed_key string, failed_source []string) {
	logger := NewLogger()

	logger.Errorw("Failed to get key-value pair",
		"failed sources", failed_source,
		"key", failed_key,
	)
}

func GetEnv(key string) string {
	var failed_source []string

	valueOS := getOSEnv(key, &failed_source)
	valueDotEnv := getDotEnv(key, &failed_source)

	if valueOS != "" {
		return valueOS
	} else if valueDotEnv != "" {
		return valueDotEnv
	}

	logger(key, failed_source)
	return ""
}

func getOSEnv(key string, failed_source *[]string) string {
	value := os.Getenv(key)
	if value == "" {
		*failed_source = append(*failed_source, "os")
		return ""
	}

	return value
}

func getDotEnv(key string, failed_source *[]string) string {
	err := godotenv.Load()
	if err != nil {
		*failed_source = append(*failed_source, ".env")
		return ""
	}

	value, ok := os.LookupEnv(key)
	if !ok {
		return ""
	}
	return value
}