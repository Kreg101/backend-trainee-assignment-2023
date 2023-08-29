package db

type User struct {
	Id             int64    `json:"id"`
	AppendSegments []string `json:"append_segments,omitempty"`
	DeleteSegments []string `json:"delete_segments,omitempty"`
}
