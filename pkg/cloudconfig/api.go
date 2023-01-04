package cloudconfig

import "os"

func GetCloudBaseURL() string {
	const CLOUD_URL = "https://cloud.mify.io"

	env := os.Getenv("MIFY_CLOUD_URL")
	if env != "" {
		return env
	}
	return CLOUD_URL
}

func GetCloudURL() string {
	return GetCloudBaseURL() + "/api"
}

func GetCloudStatsURL() string {
	return GetCloudBaseURL() + "/stats-api"
}
