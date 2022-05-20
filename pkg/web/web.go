package web

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stv0g/vand/pkg/config"
	"github.com/stv0g/vand/pkg/web/handlers"
)

const apiBase = "/api/v1"

func Run(cfg *config.Config, version, commit, date string) {
	router := gin.Default()

	router.Use(StaticMiddleware(cfg))

	router.GET(apiBase+"/config", handlers.HandleConfigWith(version, commit, date))

	server := &http.Server{
		Addr:           cfg.Web.Listen,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Listening on: http://%s", server.Addr)

	server.ListenAndServe()
}
