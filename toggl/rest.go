package toggl

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (t *Toggl) getQuestions(w http.ResponseWriter, r *http.Request) {

	questions, err := t.dbGetAllQuestions(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(questions) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Nothing found"))
		return
	}
	data, err := json.Marshal(questions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
func (t *Toggl) getQuestionsByID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("id is missing in parameters"))
		return
	}
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var question RequestQuestion
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	question.offset = ID
	questions, err := t.dbGetQuestionsByID(r.Context(), question)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Nothing found"))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(questions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
func (t *Toggl) createQuestion(w http.ResponseWriter, r *http.Request) {
	var question Question
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !isValidQuestion(question) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Question is not valid, requires atleast one correct option"))
		return
	}
	q, err := t.dbInsertQuestion(r.Context(), question)
	if err != nil {
		log.Println("[ERROR] createQuestion, error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	data, err := json.Marshal(q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (t *Toggl) updateQuestion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("id is missing in parameters"))
		return
	}

	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var question Question
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !isValidQuestion(question) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Question is not valid, requires atleast one correct option"))
		return
	}

	resp, err := t.dbUpdateQuestion(r.Context(), question, ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("[ERROR] dbUpdateQuestion, error: ", err)
		return
	}
	data, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (t *Toggl) deleteQuestion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	done, err := t.dbDeleteQuestion(r.Context(), ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("[ERROR] dbDeleteQuestion, error: ", err)
		return
	}
	if done == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Question with provided id: %d doesnt exist", ID)))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted successfuly"))
}

// isValidQuestion returns if question has atleast one correct option
func isValidQuestion(q Question) bool {
	valid := false
	for _, option := range q.Options {
		if option.Correct {
			return true
		}
	}
	return valid
}
