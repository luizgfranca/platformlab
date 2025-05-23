package server

import (
	"log"
	"net/http"
	commonmodule "synthreon/modules/common"
	configurationmodule "synthreon/modules/configuration"
	"synthreon/server/api"
	"synthreon/server/api/middleware"
	server "synthreon/server/handler"
	"synthreon/server/route"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func StartServer(addr string, configService *configurationmodule.ConfigurationService, db *gorm.DB) {
	appHandlers := &server.AppHandlers{
		ProjectAPI:        api.ProjectRESTApi(db),
		ToolAPI:           api.ToolRestAPI(db, configService),
		AuthenticationAPI: api.AuthenticationRESTApi(db, configService.AccessTokenSecret),
		WebHandler: &commonmodule.SPAHandler{
			StaticPath: configService.StaticFilesDir,
			IndexPath:  "index.html",
		},
	}

	router := mux.NewRouter()
	router.Use(middleware.GetCORSMiddleware())
	route.SetupBaseRoutes(router, appHandlers)
	route.SetupWebRoutes(router, appHandlers)

	authenticatedRouter := router.PathPrefix("/api").Subrouter()
	authenticatedRouter.Use(middleware.GetContentTypeMiddleware(middleware.ContentTypeJSON))
	authenticatedRouter.Use(middleware.GetSessionMiddleware(configService.AccessTokenSecret))
	route.SetupAuthenticatedRoutes(authenticatedRouter, appHandlers)

	log.Println("[Server] listening at", addr)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatalln("[Server] error setting up server:", err.Error())
	}
}
