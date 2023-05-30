// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stv0g/vand/pkg/store"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleWebsocketWith(store *store.Store) func(c *gin.Context) {
	return func(c *gin.Context) {
		conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("Failed to set websocket upgrade: %s", err)
			return
		}

		sups, err := store.Subscribe()
		if err != nil {
			log.Printf("Failed get updates: %+v", err)
			return
		}

		for sup := range sups {
			pl, err := json.MarshalIndent(sup, "", "  ")
			if err != nil {
				log.Printf("Failed to marshal message: %s", err)
				return
			}

			if err := conn.WriteMessage(websocket.TextMessage, pl); err != nil {
				log.Printf("Failed to write message: %s", err)
				return
			}
		}
	}
}
