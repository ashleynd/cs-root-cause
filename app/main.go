package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	_ "embed"

	_ "github.com/jackc/pgx/v4/stdlib"
)

//go:embed index.html
var index string

//go:embed script.js
var js string

type server struct {
	store  Store
	router *http.ServeMux
}

type CheckBoxItem struct {
	ID          int
	Description string
}

var todoHTML *template.Template

func main() {
	frontendPort := MustGetEnv("APP_FRONTEND_PORT")
	storeValue := GetEnv("DB", "postgres")
	databasePort := GetEnv("DB_PORT", "9091")
	databaseHost := GetEnv("DB_HOST", "localhost")
	databaseURL := fmt.Sprintf("postgres://postgres:password@%s:%s/postgres?sslmode=disable", databaseHost, databasePort)

	log.Printf("APP_FRONTEND_PORT is %s\n", frontendPort)
	log.Printf("DB_PORT is %s\n", databasePort)
	log.Printf("DB_HOST is %s\n", databaseHost)

	log.Printf("database url: %s\n", databaseURL)

	dbstore, err := createStore(storeValue, databaseURL)
	if err != nil {
		log.Fatalf("unable to open db %v", err)
	}
	defer dbstore.Close()

	srv := &server{
		store:  dbstore,
		router: http.NewServeMux(),
	}
	srv.setupRoutes()
	todoHTML = template.Must(template.New("todo").Parse(index))

	log.Printf("listening on: %s\n", frontendPort)
	if err := http.ListenAndServe(":"+frontendPort, srv.router); err != nil {
		log.Print(err)
	}
}

func (s *server) setupRoutes() {
	s.router.HandleFunc("/", s.handleIndex())
	s.router.HandleFunc("/update", s.handleUpdate())
	s.router.HandleFunc("/script.js", s.handleJS())
}

type TodoData struct {
	ScriptJS      string
	CheckBoxItems []CheckBoxItem
}

func (s *server) handleJS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/javascript")
		_, err := w.Write([]byte(js))
		if err != nil {
			logServerErr(w, err)
			return
		}
	}
}

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := s.store.List()
		if err != nil {
			log.Printf("WARN: while listing todo items: %v", err)
			items = []CheckBoxItem{}
		}

		data := TodoData{
			ScriptJS:      js,
			CheckBoxItems: items,
		}

		err = todoHTML.Execute(w, data)
		if err != nil {
			logServerErr(w, err)
		}
	}
}

func logServerErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, werr := w.Write([]byte(fmt.Sprintf("encountered err: %v", err)))
	if werr != nil {
		log.Printf("logServerErr: err writing response: %v\n", werr)
	}
	log.Printf("error:  %v \n", err)
}
func extractUpdate(r *http.Request) (*CheckBoxItem, error) {
	if r.Method != http.MethodPost {
		return nil, fmt.Errorf("extractUpdate: not a post request")
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("while reading body: %w", err)
	}

	description := string(body)
	if description == "" {
		return nil, fmt.Errorf("empty description for item")
	}

	return &CheckBoxItem{
		Description: description,
	}, nil
}

func (s *server) handleUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todoItem, err := extractUpdate(r)
		if err != nil {
			logServerErr(w, err)
			return
		}

		err = s.store.Insert(*todoItem)
		if err != nil {
			logServerErr(w, err)
			return
		}
	}
}
