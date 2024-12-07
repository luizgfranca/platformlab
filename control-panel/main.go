package main

import (
	"log"
	"net/http"
	"platformlab/controlpanel/api"
	"platformlab/controlpanel/api/middleware"
	"platformlab/controlpanel/model"
	"platformlab/controlpanel/service"
	"platformlab/controlpanel/util"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DoMigrations(db *gorm.DB) {
	db.AutoMigrate(&model.Project{})
	db.AutoMigrate(&model.Tool{})
	db.AutoMigrate(&model.User{})
}

func CreateExampleProjectsIfNotExists(db *gorm.DB) {
	s := service.Project{Db: db}

	testProjects := []model.Project{
		{Acronym: "sandbox", Name: "Sandbox", Description: "Sandbox project to test tool development."},
	}

	for i := range testProjects {
		p := testProjects[i]

		dbProject, _ := s.FindByAcronym(p.Acronym)
		if dbProject == nil {
			log.Println("saving: ", p.Acronym)
			s.Create(&p)
		}
	}
}

func CreateDefaultUserIfNotExists(db *gorm.DB, email string, password string) {
	s := service.User{Db: db}

	defaultUser, err := model.NewUser("root", email, password)
	if err != nil {
		panic(err.Error())
	}

	user, _ := s.FindByEmail(defaultUser.Email)
	if user != nil {
		return
	}

	_, err = s.Create(defaultUser)
	if err != nil {
		panic(err.Error())
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middleware", r.Method)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("trying to load application configuration")
	configuration, err := util.TryLoadApplicationConfigFromEnvironment()
	if err != nil {
		panic(err.Error())
	}

	router := mux.NewRouter()

	db, err := gorm.Open(sqlite.Open(configuration.DatabasePath), &gorm.Config{})
	if err != nil {
		panic("failed connecting to database")
	}

	println("doing database migrations")
	DoMigrations(db)
	println("done")

	println("creating example projects")
	CreateExampleProjectsIfNotExists(db)
	println("done")

	println("asseting creation of default user")
	CreateDefaultUserIfNotExists(db, configuration.RootUserEmail, configuration.RootPassword)
	println("done")

	projectAPI := api.ProjectRESTApi(db)
	toolAPI := api.ToolRestAPI(db)
	authenticationAPI := api.AuthenticationRESTApi(db, configuration.AccessTokenSecret)
	sessionMiddleware := middleware.SessionMiddleware(configuration.AccessTokenSecret)
	tableAPI := api.Table{}

	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	})

	router.HandleFunc("/api/auth/login", authenticationAPI.Login()).Methods("POST")
	router.HandleFunc("/api/tool/event", toolAPI.GetEventRresponseTEST()).Methods("POST")
	router.HandleFunc("/api/tool/provider/ws", toolAPI.ToolProviderWebsocket()).Methods("GET")
	router.HandleFunc("/api/table", tableAPI.GetTablesMetadata())

	authenticatedRouter := router.PathPrefix("/api").Subrouter()
	authenticatedRouter.Use(sessionMiddleware)
	authenticatedRouter.HandleFunc("/project", projectAPI.GetAllProjects()).Methods("GET")
	authenticatedRouter.HandleFunc("/project", projectAPI.CreateProject()).Methods("POST")
	authenticatedRouter.HandleFunc("/project/{project}/tool", projectAPI.GetToolsFromProject()).Methods("GET")
	authenticatedRouter.HandleFunc("/project/{project}/tool", projectAPI.CreateToolForProject()).Methods("POST")
	authenticatedRouter.HandleFunc("/tool/client/ws", toolAPI.ToolClientWebsocket()).Methods("GET")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/control-panel", http.StatusFound)
	})

	router.PathPrefix("/control-panel").Handler(http.StripPrefix("/control-panel", http.FileServer(http.Dir("./web/dist"))))
	router.PathPrefix("/assets").Handler(http.FileServer(http.Dir("./web/dist")))

	println("listening at :8080")
	http.ListenAndServe(":8080", corsMiddleware(router))
}
