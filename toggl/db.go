package toggl

import (
	_ "github.com/mattn/go-sqlite3"
)

func (t *Toggl) dbGetAllQuestions() []Question {
	questions := make([]Question, 0)
	rows, err := t.db.Query("SELECT id, body FROM question")
	if err != nil {
		return []Question{}
	}
	defer rows.Close()
	for rows.Next() {
		var question Question
		rows.Scan(&question.ID, &question.Body)
		questions = append(questions, question)
	}
	return questions
}

func (t *Toggl) getOptionsByQuestion(questionID int64) []Option {
	options := make([]Option, 0)
	rows, err := t.db.Query("SELECT body,correct FROM options where questionid = $1", questionID)
	if err != nil {
		return []Option{}
	}
	defer rows.Close()
	for rows.Next() {
		var option Option
		rows.Scan(&option.Body, &option.Correct)
		options = append(options, option)
	}
	return options
}
