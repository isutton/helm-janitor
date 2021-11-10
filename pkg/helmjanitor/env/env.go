package env

import (
	"os"
	"strconv"
)

func BoolOr(name string, def bool) bool {
	if name == "" {
		return def
	}
	envVal := os.Getenv(name)
	if envVal == "" {
		return def
	}
	ret, err := strconv.ParseBool(envVal)
	if err != nil {
		return def
	}
	return ret
}

func String(name string) string {
	if name == "" {
		return ""
	}
	return os.Getenv(name)
}

func StringOr(name string, def string) string {
	if name == "" {
		return def
	}
	envVal := String(name)
	if envVal == "" {
		return def
	}
	return envVal
}
