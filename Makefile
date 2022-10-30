.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

build: ### Build docker
	@docker build . -t ssv682/otus-homework3-repo
.PHONY: build

run: ### Run docker
	@docker run -p 8080:8080 ssv682/otus-homework3-repo
.PHONY: run

compose-up: ### Run docker-compose
	docker-compose up --build -d && docker-compose logs -f
.PHONY: compose-up

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

kube-up: ### Run deploy
	kubectl apply -f ./kube
.PHONY: kube-up

kube-down: ###Delere deploy
	kubectl apply -d ./kube
.PHONY: kube-down