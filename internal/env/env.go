package env

import (
	"fmt"
	"os"
)

func GetEnv(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	fmt.Print(val)
	if !ok {
		return fallback
	}
	return val

}
