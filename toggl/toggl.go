package toggl

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

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

// StartupSql exec startup sql script for the database. USED FOR TESTS
func (t *Toggl) StartupSql() error {
	query, err := os.ReadFile("../sql/startup.sql")
	if err != nil {
		return err
	}
	if _, err := t.db.Exec(string(query)); err != nil {
		return err
	}
	return nil
}

func (t *Toggl) Run() {
	log.Println("Toggl server is running on: ", t.server.Addr)
	if err := t.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Println("[ERROR] ListenAndServe failed, error: ", err)
	}
}
func (t *Toggl) Stop() {
	if err := t.db.Close(); err != nil {
		log.Println("[ERROR] Database failed to gracefuly shutdown, error:", err)
	}
	if err := t.server.Shutdown(context.Background()); err != nil {
		log.Println("[ERROR] Server failed to gracefuly shutdown, error:", err)
	}
	log.Println("Toggl terminated successfuly")

}
