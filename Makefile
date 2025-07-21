APP_NAME = food-order-api
CMD_PATH = cmd/api
SOCKET_PATH = cmd/socket
DOCKER_DEV_FILE = docker-compose.dev.yml

.PHONY: run build clean-build deps stop logs swagger test-consistency

rest:
	go run $(CMD_PATH)/main.go

socket:
	go run $(SOCKET_PATH)/main.go

build:
	go build -o bin/$(APP_NAME) $(CMD_PATH)/main.go
	go build -o bin/$(SOCKET_NAME) $(SOCKET_PATH)/main.go

clean-build:
	rm -f bin/$(APP_NAME)
	rm -f bin/$(SOCKET_NAME)
	make build

# Start all dev dependencies (MySQL, Redis, Kafka, Zookeeper)
deps:
	docker-compose -f $(DOCKER_DEV_FILE) up -d --build

stop:
	docker-compose -f $(DOCKER_DEV_FILE) down

logs:
	docker-compose -f $(DOCKER_DEV_FILE) logs -f

swagger:
	swag init --generalInfo cmd/api/main.go --output docs

# Build production images
build-prod:
	docker build -t food-order-rest --target=rest .
	docker build -t food-order-socket --target=socket .

# Start production stack
prod-up:
	docker-compose -f docker-compose.prod.yml up -d

# Stop production stack
prod-down:
	docker-compose -f docker-compose.prod.yml down -v


