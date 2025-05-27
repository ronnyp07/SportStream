package models

import "time"

type Job struct {
	TrackInfo TrackInfo `json:"track_info"`
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	LastRun   time.Time `json:"last_run"`
	NextRun   time.Time `json:"next_run"`
	Tags      []string  `json:"tags"`
}

type TrackInfo struct {
	ID        string    `json:"id"`
	Source    string    `json:"source"`
	Timestamp time.Time `json:"timestamp"`
}

type TrackingLog struct {
	TrackID  *string `json:"track_id,omitempty"`
	Hostname string  `json:"hostname"`
	Type     *string `json:"type,omitempty"`
	EventID  *int    `json:"event_id,omitempty"`
	ID       *int    `json:"id,omitempty"`
	Message  string  `json:"message"`
	File     string  `json:"file"`
	Line     int     `json:"line"`
}
