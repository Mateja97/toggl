package toggl

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (t *Toggl) getQuestions(w http.ResponseWriter, r *http.Request) {
	var resp []responseQuestion
	questions := t.dbGetAllQuestions()
	if len(questions) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Nothing found"))
		return
	}
	for _, q := range questions {
		var question responseQuestion
		question.Body = q.Body
		options := t.getOptionsByQuestion(q.ID)
		if len(options) == 0 {
			continue
		}
		resp = append(resp, question)
	}
}

func (t *Toggl) createQuestion(w http.ResponseWriter, r *http.Request) {
	var question model.Question
	json.NewDecoder(r.Body).Decode(&question)
	result, err := t.db.Exec("INSERT INTO question (body) VALUES (?)", question.Body)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	question.ID = id
	for _, option := range question.Options {
		_, err := t.db.Exec("INSERT INTO option (question, body, correct) VALUES (?,?,?)", id, option.Body, option.Correct)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
	}
	json.NewEncoder(w).Encode(question)
}

func (t *Toggl) updateQuestion(w http.ResponseWriter, r *http.Request) {
	var question Question
	json.NewDecoder(r.Body).Decode(&question)
	_, err := t.db.Exec("UPDATE question SET body = ? WHERE id = ?", question.Body, question.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	_, err = t.db.Exec("DELETE FROM option WHERE question = ?", question.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	for _, option := range question.Options {
		_, err := t.db.Exec("INSERT INTO option (question, body, correct) VALUES (?,?,?)", question.ID, option.Body, option.Correct)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
	}
	json.NewEncoder(w).Encode(question)
}

func (t *Toggl) deleteQuestion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	_, err := t.db.Exec("DELETE FROM question WHERE id = ?", id)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	_, err = t.db.Exec("DELETE FROM option WHERE question = ?", id)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"result": "success"})
}
