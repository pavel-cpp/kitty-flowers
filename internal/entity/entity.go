package entity

import (
	"time"

	"github.com/google/uuid"
)

type (
	User struct {
		ID       uuid.UUID
		ChatID   int
		Username string
	}

	UserNotification struct {
		User
		NotificationID int
		CurrentRun     time.Time
	}
)
