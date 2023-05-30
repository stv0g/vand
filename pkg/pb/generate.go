// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package pb

//go:generate protoc --go_out=. --go_opt=paths=source_relative battery.proto obd2.proto environment.proto gps.proto modem.proto solar.proto vand.proto
