run: ## build/run flight_tracker
	docker build -t flight_tracker .
	docker run -p 8080:8080 flight_tracker

stop: ## shutdown & remove flight_tracker
	docker stop flight_tracker
	docker rm flight_tracker

test: ## test flight_tracker
	go test -v -timeout 30s ./...

coverage: ## test coverage flight_tracker
	go test -timeout 30s -coverprofile=go-code-cover ./...

