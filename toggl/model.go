package toggl

type Option struct {
	Body    string `json:"body"`
	Correct bool   `json:"correct"`
}

// Same struct as in the database
type Question struct {
	ID      int64    `json:"id"`
	Body    string   `json:"body"`
	Options []Option `json:"options"`
}
type RequestQuestion struct {
	offset int64
	Limit  int64 `json:"limit"`
}
