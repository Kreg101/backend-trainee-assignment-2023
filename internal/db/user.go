package db

type User struct {
	Id          int64    `json:"id"`
	NewSegments []string `json:"new_segments,omitempty"`
	OldSegments []string `json:"old_segments,omitempty"`
}
