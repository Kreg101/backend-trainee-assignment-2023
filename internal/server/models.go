package server

// User represents id and segments for adding, deleting, etc.
type User struct {
	Id         int64    `json:"id"`
	Segments   []string `json:"segments,omitempty"`
	ActiveTime int64    `json:"active_time,omitempty"`
	Year       int      `json:"year,omitempty"`
	Month      int      `json:"month,omitempty"`
}

// Msg represents description of problem in request
type Msg struct {
	Text string `json:"text"`
}

// Segment structure for json unmarshalling
type Segment struct {
	Name string `json:"segment"`
}

// TimeUser for history
type TimeUser struct {
	Id          int64
	SegmentName string
	TimeIn      string
	TimeOut     string
}
