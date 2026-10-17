package config

import (
	"fmt"
	"os"
)

type Config struct {
	ApiPort            string
	DatabaseURL        string
	DownloadDir        string
	CookiesFromBrowser string
	CookiesFile        string
}

func LoadConfig() (Config, error) {
	apiPort, err := required("API_PORT")
	if err != nil {
		return Config{}, err
	}

	dbHost, err := required("DB_HOST")
	if err != nil {
		return Config{}, err
	}

	dbPort, err := required("DB_PORT")
	if err != nil {
		return Config{}, err
	}

	dbUser, err := required("DB_USER")
	if err != nil {
		return Config{}, err
	}

	dbPassword, err := required("DB_PASSWORD")
	if err != nil {
		return Config{}, err
	}

	dbName, err := required("DB_NAME")
	if err != nil {
		return Config{}, err
	}

	downloadDir, err := required("DOWNLOAD_DIR")
	if err != nil {
		return Config{}, err
	}

	dbURL := buildDatabaseURL(dbHost, dbPort, dbUser, dbPassword, dbName)

	cookiesFromBrowser := os.Getenv("YT_COOKIES_FROM_BROWSER")
	cookiesFile := os.Getenv("YT_COOKIES_FILE")

	return Config{
		ApiPort:            apiPort,
		DatabaseURL:        dbURL,
		DownloadDir:        downloadDir,
		CookiesFromBrowser: cookiesFromBrowser,
		CookiesFile:        cookiesFile,
	}, nil
}

func required(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable %s is required", key)
	}
	return value, nil
}

func buildDatabaseURL(host, port, user, password, name string) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, name,
	)
}
