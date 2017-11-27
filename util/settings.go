package util

import (
	"os"
)

func GetAPIURL() string {
	if os.Getenv("API_URL") != "" {
		return os.Getenv("API_URL")
	}
	return "http://api.tourney.aback.es:9090/api/v2"
}

func GetStorageLocation() string {
	if os.Getenv("TOURNEY_STORAGE_LOCATION") != "" {
		return os.Getenv("TOURNEY_STORAGE_LOCATION")
	}
	return "tourney_storage"
}
