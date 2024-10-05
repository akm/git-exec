package main

import (
	"os"
	"strings"
)

func getEnvBool(key string) bool {
	v := strings.ToLower(os.Getenv(key))
	switch v {
	case "true", "1", "yes", "on":
		return true
	default:
		return false
	}
}
