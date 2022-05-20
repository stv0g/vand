TARGET = vand
TARGET_HOST ?= root@rpi.bus.0l.de

REMOTE = ssh -t -i ~/.ssh/id_rsa $(TARGET_HOST)

export GOOS ?= linux
export GOARCH ?= arm

ifeq ($(GOARCH),arm)
  export GOARM ?= 7
endif

all: run

frontend-code:
	npm --prefix frontend/ run-script build

build: frontend-code
	go build -tags embed_frontend -o $(TARGET) ./cmd/

deploy: build
	rsync --progress $(TARGET) $(TARGET_HOST):/usr/local/bin/$(TARGET)
	rsync --progress etc/vand.yaml $(TARGET_HOST):/etc/

run: deploy
	$(REMOTE) /usr/local/bin/$(TARGET) gps

.PHONY: build run
