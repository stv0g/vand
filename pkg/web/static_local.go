// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

//go:build !embed_frontend

package web

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/stv0g/vand/pkg/config"
)

// StaticMiddleware serves static assets package by Webpack
func StaticMiddleware(cfg *config.Config) gin.HandlerFunc {
	return static.Serve("/", static.LocalFile(cfg.Web.Static, false))
}
