package cloudconfig

import "os"

const DEFAULT_CLOUD_URL = "https://cloud.mify.io"
const DEFAULT_CLOUD_API_URL = "https://cloud.mify.io/api"
const DEFAULT_STATS_API_URL = "https://cloud.mify.io/api-stats"

func GetCloudUrl() string {
	return DEFAULT_CLOUD_URL
}

func GetCloudApiURL() string {
	env := os.Getenv("MIFY_CLOUD_API_URL")
	if env != "" {
		return env
	}
	return DEFAULT_CLOUD_API_URL
}

func GetStatsApiUrl() string {
	env := os.Getenv("MIFY_STATS_API_URL")
	if env != "" {
		return env
	}
	return DEFAULT_STATS_API_URL
}
