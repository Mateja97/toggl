package toggl

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Toggl struct {
	db     *sql.DB
	server *http.Server
}

func (t *Toggl) Init(port string) {
	var err error
	t.db, err = sql.Open("sqlite3", "./question.db")
	if err != nil {
		panic(err)
	}
	defer t.db.Close()

	router := mux.NewRouter()

	// JWT Middleware
	router.Use(jwtMiddleware)

	// Endpoints
	router.HandleFunc("/questions", t.getQuestions).Methods("GET")
	router.HandleFunc("/questions", t.createQuestion).Methods("POST")
	router.HandleFunc("/questions/{id}", t.updateQuestion).Methods("PUT")
	router.HandleFunc("/questions/{id}", t.deleteQuestion).Methods("DELETE")
	t.server = &http.Server{
		Handler: router,
		Addr:    port,
	}

}

func (t *Toggl) Run() {
	if err := t.server.ListenAndServe(); err != nil {
		fmt.Errorf("ListenAndServe failed, error:", err)
	}
	fmt.Println("Toggl server is running successfuly")
}
func (t *Toggl) Stop() {
	if err := t.db.Close(); err != nil {
		fmt.Errorf("Database failed to gracefuly shutdown, error:", err)
	}
	if err := t.server.Shutdown(context.Background()); err != nil {
		fmt.Errorf("Http server failed to gracefuly shutdown, error:", err)
	}
	fmt.Println("Toggl terminated successfuly")

}
