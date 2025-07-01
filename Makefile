APP_NAME = food-order-api
CMD_PATH = cmd/api
DOCKER_DEV_FILE = docker-compose.dev.yml

.PHONY: run build clean-build deps stop logs

run:
	go run $(CMD_PATH)/main.go

build:
	go build -o bin/$(APP_NAME) $(CMD_PATH)/main.go

clean-build:
	rm -f bin/$(APP_NAME)
	make build

deps:
	docker-compose -f $(DOCKER_DEV_FILE) up -d

stop:
	docker-compose -f $(DOCKER_DEV_FILE) down

logs:
	docker-compose -f $(DOCKER_DEV_FILE) logs -f
