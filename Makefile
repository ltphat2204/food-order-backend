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

deps:
	docker-compose -f $(DOCKER_DEV_FILE) up -d

stop:
	docker-compose -f $(DOCKER_DEV_FILE) down

logs:
	docker-compose -f $(DOCKER_DEV_FILE) logs -f

swagger:
	swag init --generalInfo cmd/api/main.go --output docs

# Test aggregate consistency
test-consistency:
	go test ./internal/app/order/service/ -v -run="TestAggregate"

