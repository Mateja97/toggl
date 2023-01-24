package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/Mateja97/toggl/toggl"
	"github.com/golang-jwt/jwt/v4"
)

var CreateQuestion = toggl.Question{
	Body: "Who is serbian president?",
	Options: []toggl.Option{
		{
			Body:    "Novak Djokovic",
			Correct: true,
		}, {
			Body:    "Natasa Bekvalac",
			Correct: false,
		}, {
			Body:    "Aleksandar Vucic",
			Correct: false,
		},
	},
}
var CreateQuestion2 = toggl.Question{
	Body: "What is the capitol city of Serbia?",
	Options: []toggl.Option{
		{
			Body:    "Belgrade",
			Correct: true,
		}, {
			Body:    "Kragujevac",
			Correct: false,
		}, {
			Body:    "Nis",
			Correct: false,
		},
	},
}
var tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.Yt3L4ocGM7nuYpi0jruAKsPvMkns1iQqGqJPgbYRTVk"
var port = "8081"
var questionsURL = fmt.Sprintf("http://localhost%s/questions/", port)
var PaginationQuestion = &toggl.RequestQuestion{
	Limit: 1,
}

func TestCreateQuestion(t *testing.T) {
	//Init service
	toggleApp := toggl.Toggl{}
	if err := toggleApp.Init(port); err != nil {
		t.Errorf("[ERROR] Service failed to init, error: %s", err)
	}
	err := toggleApp.StartupSql()
	if err != nil {
		t.Errorf("Error on startup sql, %s", err)
	}
	serviceRunning := make(chan struct{})
	serviceDone := make(chan struct{})
	go func() {
		close(serviceRunning)
		toggleApp.Run()
		defer close(serviceDone)
	}()
	<-serviceRunning
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("toggl"), nil
	})
	if err != nil || !token.Valid {
		t.Errorf("Not Authorized")
	}
	client := &http.Client{}
	//CREATE QUESTION 1
	///
	///
	questionJSON, _ := json.Marshal(CreateQuestion)
	req, err := http.NewRequest("POST", questionsURL, bytes.NewBuffer(questionJSON))
	if err != nil {
		t.Errorf("Got error %s", err.Error())
	}
	req.Header.Set("Authorization", tokenString)
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("Error sending POST request: %s", err)
	}
	//CREATE QUESTION 2
	///
	///
	question2JSON, _ := json.Marshal(CreateQuestion2)
	req, err = http.NewRequest("POST", questionsURL, bytes.NewBuffer(question2JSON))
	if err != nil {
		t.Errorf("Got error %s", err.Error())
	}
	req.Header.Set("Authorization", tokenString)
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error sending POST request: %s", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var question toggl.Question
	_ = json.Unmarshal(body, &question)

	//CHECK RESULT
	///
	///
	paginationJSON, _ := json.Marshal(PaginationQuestion)
	url := fmt.Sprintf("%s%d", questionsURL, question.ID+1)
	req2, err := http.NewRequest("GET", url, bytes.NewBuffer(paginationJSON))
	if err != nil {
		t.Errorf("Got error %s", err.Error())
	}
	req2.Header.Set("Authorization", tokenString)
	resp2, err := client.Do(req2)
	if err != nil {
		t.Errorf("Error sending GET request: %s", err)
	}
	body, _ = ioutil.ReadAll(resp2.Body)
	var returned []toggl.Question
	var returnedQuestion toggl.Question
	_ = json.Unmarshal(body, &returned)
	returnedQuestion = returned[0]
	if returnedQuestion.ID != question.ID {
		t.Errorf("Wrong id: %d, should be %d", returnedQuestion.ID, question.ID)
	}
	toggleApp.Stop()
	<-serviceDone
}
func TestUpdateQuestion(t *testing.T) {
	//Initservice
	toggleApp := toggl.Toggl{}
	if err := toggleApp.Init(port); err != nil {
		t.Errorf("[ERROR] Service failed to init, error: %s", err)
	}
	err := toggleApp.StartupSql()
	if err != nil {
		t.Errorf("Error on startup sql, %s", err)
	}
	serviceRunning := make(chan struct{})
	serviceDone := make(chan struct{})
	go func() {
		close(serviceRunning)
		toggleApp.Run()
		defer close(serviceDone)
	}()
	<-serviceRunning
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("toggl"), nil
	})
	if err != nil || !token.Valid {
		t.Errorf("Not Authorized")
	}
	client := &http.Client{}
	//CREATE QUESTION 1
	///
	///
	questionJSON, _ := json.Marshal(CreateQuestion)
	req, err := http.NewRequest("POST", questionsURL, bytes.NewBuffer(questionJSON))
	if err != nil {
		t.Errorf("Got error %s", err.Error())
	}
	req.Header.Set("Authorization", tokenString)
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error sending POST request: %s", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var question toggl.Question
	_ = json.Unmarshal(body, &question)
	//UPDATE QUESTION 1
	///
	///
	uquestionJSON, _ := json.Marshal(CreateQuestion2)
	url := fmt.Sprintf("%s%d", questionsURL, question.ID)
	req2, err := http.NewRequest("PUT", url, bytes.NewBuffer(uquestionJSON))
	if err != nil {
		t.Errorf("Got error %s", err.Error())
	}
	req2.Header.Set("Authorization", tokenString)
	resp2, err := client.Do(req2)
	if err != nil {
		t.Errorf("Error sending PUT request: %s", err)
	}
	body, _ = ioutil.ReadAll(resp2.Body)
	var returnedQuestion toggl.Question
	_ = json.Unmarshal(body, &returnedQuestion)

	if returnedQuestion.ID != question.ID {
		t.Errorf("Wrong id updated")
	}
	//CHECK RESULTS
	///
	///
	paginationJSON, _ := json.Marshal(PaginationQuestion)
	url = fmt.Sprintf("%s%d", questionsURL, question.ID+1)
	req3, err := http.NewRequest("GET", url, bytes.NewBuffer(paginationJSON))
	if err != nil {
		t.Errorf("Got error %s", err.Error())
	}
	req3.Header.Set("Authorization", tokenString)
	resp3, err := client.Do(req3)
	if err != nil {
		t.Errorf("Error sending GET request: %s", err)
	}
	body, _ = ioutil.ReadAll(resp3.Body)
	var returned []toggl.Question
	_ = json.Unmarshal(body, &returned)
	if returned[0].ID != question.ID {
		t.Errorf("Wrong id: %d, should be %d", returnedQuestion.ID, question.ID)
	}
	if returned[0].Body != returnedQuestion.Body || !reflect.DeepEqual(returned[0].Options, returnedQuestion.Options) {
		t.Errorf("Question not updated got: %v, should be %v", returned[0], returnedQuestion)
	}
}
func TestDeleteQuestion(t *testing.T) {
	//Init service
	toggleApp := toggl.Toggl{}
	if err := toggleApp.Init(port); err != nil {
		t.Errorf("[ERROR] Service failed to init, error: %s", err)
	}
	err := toggleApp.StartupSql()
	if err != nil {
		t.Errorf("Error on startup sql, %s", err)
	}
	serviceRunning := make(chan struct{})
	serviceDone := make(chan struct{})
	go func() {
		close(serviceRunning)
		toggleApp.Run()
		defer close(serviceDone)
	}()
	<-serviceRunning
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("toggl"), nil
	})
	if err != nil || !token.Valid {
		t.Errorf("Not Authorized")
	}
	client := &http.Client{}
	//CREATE QUESTION 1
	///
	///
	questionJSON, _ := json.Marshal(CreateQuestion)
	req, err := http.NewRequest("POST", questionsURL, bytes.NewBuffer(questionJSON))
	if err != nil {
		t.Errorf("Got error %s", err.Error())
	}
	req.Header.Set("Authorization", tokenString)
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error sending POST request: %s", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var question toggl.Question
	_ = json.Unmarshal(body, &question)
	//DELETE QUESTION 1
	///
	///
	url := fmt.Sprintf("%s%d", questionsURL, question.ID)
	req2, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Errorf("Got error %s", err.Error())
	}
	req2.Header.Set("Authorization", tokenString)
	_, err = client.Do(req2)
	if err != nil {
		t.Errorf("Error sending DELETE request: %s", err)
	}
	//CHECK RESULT
	///
	///
	paginationJSON, _ := json.Marshal(PaginationQuestion)
	url = fmt.Sprintf("%s%d", questionsURL, question.ID+1)
	req3, err := http.NewRequest("GET", url, bytes.NewBuffer(paginationJSON))
	if err != nil {
		t.Errorf("Got error %s", err.Error())
	}
	req3.Header.Set("Authorization", tokenString)
	resp3, err := client.Do(req3)
	if err != nil {
		t.Errorf("Error sending GET request: %s", err)
	}
	body, _ = ioutil.ReadAll(resp3.Body)
	var returned []toggl.Question
	_ = json.Unmarshal(body, &returned)

	if len(returned) != 0 && returned[0].ID == question.ID {
		t.Errorf("question not deleted, got same question back: %v", returned[0])
	}
}
