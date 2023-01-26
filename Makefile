.PHONY: help build run compose-up compose-down kube-up kube-down migration

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build: ### Build docker
	@docker build -f ./user-service/Dockerfile ./ -t ssv682/user-service && build ./auth-service -t ssv682/auth-service

run: ### Run docker
	@docker run -p 8080:8080 ssv682/user-service

compose-up: ### Run docker-compose
	docker-compose up --build -d && docker-compose logs -f

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans

kube-up: ### Run deploy
	kubectl apply -f ./kube

kube-down: ### Delete deploy
	kubectl apply -d ./kube

migration: ### Create next migration
	@migrate create -ext sql -dir migrations/db -format "2006010215" $(name)
