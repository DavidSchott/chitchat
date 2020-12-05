package data

import "time"

type Client struct {
	Username     string    `json:"username"`
	Color        string    `json:"color"` // Not being used due to possible connection failure...
	LastActivity time.Time `json:"last_activity"`
}
