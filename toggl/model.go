package toggl

type Option struct {
	ID         int64  `json:"id"`
	QuestionID int64  `json:"questionid"`
	Body       string `json:"body"`
	Correct    bool   `json:"correct"`
}

type Question struct {
	ID   int64  `json:"id"`
	Body string `json:"body"`
}
type responseOption struct {
	Body    string `json:"body"`
	Correct bool   `json:"correct"`
}
type responseQuestion struct {
	Body    string           `json:"body"`
	Options []responseOption `json:"options"`
}
