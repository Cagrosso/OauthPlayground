package server

import "time"

type sessionTracker struct {
	AccessToken string
	TimeOut     time.Time
}
