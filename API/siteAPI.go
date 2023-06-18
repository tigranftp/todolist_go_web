package API

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"site/config"
	"site/db"
)

const DefPORT = ":8080"

type API struct {
	router *mux.Router
	config *config.Config
	store  *db.Store
}

func NewSiteAPI() (*API, error) {
	res := new(API)
	res.router = mux.NewRouter()
	var err error
	res.config, err = config.GetConfig()
	return res, err
}

func (a *API) Start() error {
	a.configureRouter()
	a.configureDB()
	if err := a.store.Open(); err != nil {
		return err
	}
	return http.ListenAndServe(DefPORT, a.router)
}

func (a *API) Stop() {
	fmt.Println("Stopping API...")
	a.store.Close()
	fmt.Println("API stopped...")
}

func (a *API) configureDB() {
	a.store = db.New(a.config)
}

func (a *API) configureRouter() {
	a.router.Use(CORS)
	a.router.HandleFunc("/sign_in", a.handleSignIn())
	a.router.HandleFunc("/test_parse_token", a.testParseToken())
	a.router.HandleFunc("/sign_up", a.handleSignUp())
	a.router.HandleFunc("/get_username_by_token", a.handleGetUsernameByToken())
	a.router.HandleFunc("/add_task_for_user", a.handleAddTaskForUser())
	a.router.HandleFunc("/get_tasks_of_user", a.handleGetTasksOfUser())
	a.router.HandleFunc("/update_todo_list_item", a.handleUpdateTodoListItem())
	a.router.HandleFunc("/delete_task_by_id", a.handleDeleteTaskByID())
	a.router.HandleFunc("/sign_in_by_refresh", a.handleSignInByRefresh())
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", `http://localhost:63342`)
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Next
		next.ServeHTTP(w, r)
		return
	})
}
