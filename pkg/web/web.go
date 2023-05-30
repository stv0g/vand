// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package web

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stv0g/vand/pkg/config"
	"github.com/stv0g/vand/pkg/store"
	"github.com/stv0g/vand/pkg/web/handlers"
)

const apiBase = "/api/v1"

func Run(cfg *config.Config, store *store.Store, version, commit, date string) error {
	router := gin.Default()

	router.Use(StaticMiddleware(cfg))

	router.GET(apiBase+"/config", handlers.HandleConfigWith(version, commit, date))
	router.GET(apiBase+"/state", handlers.HandleStateWith(store))
	router.GET(apiBase+"/ws", handlers.HandleWebsocketWith(store))

	server := &http.Server{
		Addr:           cfg.Web.Listen,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Listening on: http://%s", server.Addr)

	return server.ListenAndServe()
}
