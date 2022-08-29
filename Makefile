run: ## build/run flight_tracker
	go run ./cmd/server.go

docker_run: ## build/run flight_tracker as docker  image
	docker build -t flight_tracker .
	docker run -p 8080:8080 flight_tracker

docker_stop: ## shutdown & remove flight_tracker docker image
	docker stop flight_tracker
	docker rm flight_tracker

test: ## test flight_tracker
	go test -v -timeout 30s ./...

coverage: ## test coverage flight_tracker
	go test -timeout 30s -coverprofile=go-code-cover ./...

