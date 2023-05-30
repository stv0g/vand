// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

//go:build embed_frontend

package frontend

import (
	"embed"
)

//go:embed build
var Files embed.FS
