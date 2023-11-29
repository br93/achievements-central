up: build_network up_db
	@echo "Starting service compose..."
	docker-compose -f ./project/service-compose.yml up -d
	@echo "Done"

up_build: down build_auth
	@echo "Building and starting..."
	docker-compose -f ./project/service-compose.yml up --build -d
	@echo "Done"

up_db:
	@echo "Starting db compose..."
	docker-compose -f ./project/db-compose.yml up -d
	@echo "Migrating db..."
	docker-compose -f ./project/db-compose.yml --profile tools run --rm migrate up

down:
	@echo "Stopping compose..."
	docker-compose -f ./project/db-compose.yml down
	docker-compose -f ./project/service-compose.yml down
	@echo "Done"

build_auth:
	@echo "Building auth..."
	cd ./auth-service && env GOOS=linux CGO_ENABLED=0 go build -o app ./cmd/api
	@echo "Done"

build_network:
	@echo "Creating network..."
	docker network create -d bridge achievements-central-network || true
	@echo "Done"

build_clean:
	@echo "Cleaning untag images..."
	docker rmi `docker images | grep "<none>" | awk {'print $3'}`
	@echo "Done"

all: down up_build up