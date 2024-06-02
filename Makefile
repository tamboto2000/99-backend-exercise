prepare-configs:
	cp ./deployment/docker-compose.example.env ./deployment/docker-compose.env
	cp ./configs/user_svc.example.yaml ./configs/user_svc.yaml
	cp ./configs/public_api_svc.example.yaml ./configs/public_api_svc.yaml
	cp ./configs/migration.example.yaml ./configs/migration.yaml

deploy:
	@mkdir -p ./deployment/_volumes/listing-svc
	@touch ./deployment/_volumes/listing-svc/listings.db
	@docker compose --env-file ./deployment/docker-compose.env -f ./deployment/docker-compose.yaml -p "99-backend-service" up --build -d

shutdown:
	@docker compose --env-file ./deployment/docker-compose.env -f ./deployment/docker-compose.yaml -p "99-backend-service" down

build-user-svc:
	@echo "building binary..."
	@go build -o ./bin/user-svc ./cmd/user/main.go
	@echo "finished"
	
run-user-svc: build-user-svc
	@./bin/user-svc

build-public-api-svc:
	@echo "building binary..."
	@go build -o ./bin/public-api-svc ./cmd/public_api/main.go
	@echo "finished"

run-public-api-svc: build-public-api-svc
	@./bin/public-api-svc