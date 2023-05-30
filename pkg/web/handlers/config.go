// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package handlers

import (
	"github.com/gin-gonic/gin"
)

type respBuild struct {
	Version string `json:"version"`
	Date    string `json:"date"`
	Commit  string `json:"commit"`
}

type configResponse struct {
	Build respBuild `json:"build"`
}

// HandleConfigWith returns runtime configuration to the frontend
func HandleConfigWith(version, commit, date string) func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(200, &configResponse{
			Build: respBuild{
				Commit:  commit,
				Version: version,
				Date:    date,
			},
		})
	}
}
