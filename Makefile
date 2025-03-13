# SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
# SPDX-License-Identifier: Apache-2.0

TARGET = vand
TARGET_HOST ?= root@rpi.bus.0l.de

REMOTE = ssh -t -i ~/.ssh/id_rsa $(TARGET_HOST)

export CGO_ENABLED ?= 1

all: backend

frontend:
	npm --prefix frontend/ install
	npm --prefix frontend/ run-script build

backend:
	go build -tags embed_frontend,virtual -o $(TARGET) ./cmd/

deploy: backend
	rsync --progress $(TARGET) $(TARGET_HOST):/usr/local/bin/$(TARGET)
	rsync --progress etc/vand.yaml $(TARGET_HOST):/etc/

restart: deploy
	$(REMOTE) systemctl restart vand@gps

run: deploy
	$(REMOTE) /usr/local/bin/$(TARGET) gps

.PHONY: backend frontend run deploy restart all
