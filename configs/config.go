package configs

import (
	"os"
	"strings"
	"time"
)

type (
	Database struct {
		Driver   string
		DBString string
	}

	PollenAI struct {
		Token         string
		DefaultModel  string
		GenerateModel string
	}

	TG struct {
		Token string
	}

	Notifications struct {
		NotificationTimes []time.Time
		CheckFrequency    string
	}

	Config struct {
		DB            Database
		PollenAI      PollenAI
		TG            TG
		Notifications Notifications
	}
)

func LoadNotifications() Notifications {
	rawTimes := os.Getenv("DEFAULT_NOTIFY_TIMES")
	timesStrs := strings.Split(rawTimes, ",")
	if len(timesStrs) == 0 {
		panic("No times in DEFAULT_NOTIFY_TIMES")
	}
	times := make([]time.Time, 0, len(timesStrs))
	for _, timeStr := range timesStrs {
		time, err := time.Parse(time.TimeOnly, timeStr)
		if err != nil {
			panic("Incorrect time format")
		}
		times = append(times, time)
	}

	return Notifications{
		NotificationTimes: times,
		CheckFrequency:    os.Getenv("CHECK_TASKS_FREQUENCY"),
	}
}

func LoadPollenAI() PollenAI {
	defaultModel := os.Getenv("POLLENAI_DEFAULT_MODEL")
	if defaultModel == "" {
		panic("DEFAULT_MODEL environment variable not set")
	}
	generateModel := os.Getenv("POLLENAI_GENERATE_MODEL")
	if generateModel == "" {
		generateModel = defaultModel
	}
	return PollenAI{
		Token:         os.Getenv("POLLENAI_TOKEN"),
		DefaultModel:  generateModel,
		GenerateModel: generateModel,
	}
}

func LoadDB() Database {
	return Database{
		os.Getenv("GOOSE_DRIVER"),
		os.Getenv("GOOSE_DBSTRING"),
	}
}

func LoadConfig() Config {
	return Config{
		DB:            LoadDB(),
		PollenAI:      LoadPollenAI(),
		TG:            TG{Token: os.Getenv("TG_TOKEN")},
		Notifications: LoadNotifications(),
	}
}
