package toggl

import (
	"context"
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Toggl struct {
	db     *sql.DB
	server *http.Server
}

var togglDB = "toggl.db"

func (t *Toggl) Init(port string) error {
	var err error
	t.db, err = sql.Open("sqlite3", togglDB)
	if err != nil {
		return err
	}
	router := mux.NewRouter()

	// JWT Middleware
	router.Use(jwtMiddleware)

	// Endpoints
	router.HandleFunc("/questions", t.getQuestions).Methods("GET")
	router.HandleFunc("/questions/{id}", t.getQuestionsByID).Methods("GET")
	router.HandleFunc("/questions/", t.createQuestion).Methods("POST")
	router.HandleFunc("/questions/{id}", t.updateQuestion).Methods("PUT")
	router.HandleFunc("/questions/{id}", t.deleteQuestion).Methods("DELETE")
	t.server = &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	return nil
}

func (t *Toggl) StartupSql() error {
	query, err := ioutil.ReadFile("../sql/startup.sql")
	if err != nil {
		return err
	}
	if _, err := t.db.Exec(string(query)); err != nil {
		return err
	}
	return nil
}

func (t *Toggl) Run() error {
	log.Println("Toggl server is running on: ", t.server.Addr)
	if err := t.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Println("[ERROR] ListenAndServe failed, error: ", err)
		return err
	}
	return nil
}
func (t *Toggl) Stop() {
	if err := t.db.Close(); err != nil {
		log.Println("[ERROR] Database failed to gracefuly shutdown, error:", err)
	}
	if err := t.server.Shutdown(context.Background()); err != nil {
		log.Println("[ERROR] Http server failed to gracefuly shutdown, error:", err)
	}
	log.Println("Toggl terminated successfuly")

}
