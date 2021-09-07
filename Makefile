.DEFAULT_GOAL := build

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build:  ## Build docker eiladin/go-simple-startpage:$(tag)  with DockerBuildKet
	DOCKER_BUILDKIT=1 docker build -f experimental.Dockerfile --build-arg version=$(tag) . -t eiladin/go-simple-startpage:$(tag) 

publish: build ## Push docker eiladin/go-simple-startpage:$(tag) to docker hub
	docker push eiladin/go-simple-startpage:$(tag)
  
classicBuild: ## Build docker eiladin/go-simple-startpage:$(tag)
	docker build -f Dockerfile  --build-arg version=$(tag) . -t eiladin/go-simple-startpage:$(tag)

classicPublish: classicBuild ## Push docker eiladin/go-simple-startpage:$(tag) to docker hub with a classic docker build
	docker push eiladin/go-simple-startpage:$(tag)

test: ## Run unit tests
	go test ./... -version

test-coverage: ## Gather code coverage
	go test ./... -coverprofile=coverage.out -covermode=atomic