package main

import (
	"fmt"
	"os"
)

// GetEnv returns the value of the environment variable
// or the default value if it is empty.
func GetEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// MustGetEnv returns the value of the environment variable
// if it is not empty, or panics otherwise.
func MustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("env not set: %s", key))
	}
	return v
}
