package env

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func GetStringEnv(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val

}

func GetIntEnv(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return intVal

}

func GetTimeEnv(key string, fallback time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	fmt.Print(val)
	if !ok {
		return fallback
	}
	duration, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}
	return duration

}
