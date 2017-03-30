package config

import (
	"os"
	"strings"
	"time"

	"log"

	"github.com/joho/godotenv"
	"gopkg.in/mgo.v2"
)

// Configuration for whole application
type Configuration struct {
	Google *GoogleConfiguration
	Slack  *SlackConfiguration
	Github *GithubConfiguration
	Mongo  *mgo.DialInfo
	App    *ApplicationConfiguration
}

// GoogleConfiguration settings
type GoogleConfiguration struct {
	ClientID     string
	ClientSecret string
}

// SlackConfiguration settings
type SlackConfiguration struct {
	ClientID     string
	ClientSecret string
}

//GithubConfiguration settings
type GithubConfiguration struct {
	ClientID     string
	ClientSecret string
}

// ApplicationConfiguration for server
type ApplicationConfiguration struct {
	URL string
}

// Settings is Configuration instance
var Settings *Configuration

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print(err.Error())
	}

	appURL := os.Getenv("APP_URL")
	port := os.Getenv("PORT")

	if port != "" {
		appURL = ":" + port
	}

	Settings = &Configuration{
		Mongo: &mgo.DialInfo{
			Addrs:    strings.Split(os.Getenv("MONGO_HOSTS"), ","),
			Timeout:  10 * time.Second,
			Database: os.Getenv("MONGO_DATABASE"),
			Username: os.Getenv("MONGO_USERNAME"),
			Password: os.Getenv("MONGO_PASSWORD"),
		},
		Google: &GoogleConfiguration{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		},
		Slack: &SlackConfiguration{
			ClientID:     os.Getenv("SLACK_CLIENT_ID"),
			ClientSecret: os.Getenv("SLACK_CLIENT_SECRET"),
		},
		Github: &GithubConfiguration{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		},
		App: &ApplicationConfiguration{
			URL: appURL,
		},
	}
}
