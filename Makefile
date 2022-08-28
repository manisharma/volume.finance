run: ## build/run flight_tracker
	docker build -t flight_tracker .
	docker run -p 8080:8080 flight_tracker

stop: ## shutdown & remove flight_tracker
	docker stop flight_tracker
	docker rm flight_tracker