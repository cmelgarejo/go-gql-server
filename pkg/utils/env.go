package utils

import (
	"os"
	"strconv"

	"github.com/cmelgarejo/go-gql-server/internal/logger"
)

// MustGet will return the env or panic if it is not present
func MustGet(k string) string {
	v := os.Getenv(k)
	if v == "" {
		logger.MissingArg(k)
		logger.Panic("ENV missing, key: " + k)
	}
	return v
}

// MustGetBool will return the env as boolean or panic if it is not present
func MustGetBool(k string) bool {
	v := os.Getenv(k)
	if v == "" {
		logger.MissingArg(k)
		logger.Panic("ENV missing, key: " + k)
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		logger.MissingArg(k)
		logger.Panic("ENV err: [" + k + "]" + err.Error())
	}
	return b
}
