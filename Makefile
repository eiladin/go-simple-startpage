.DEFAULT_GOAL := build

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build:  ## Build docker eiladin/go-simple-startpage:$(tag)  with DockerBuildKet
	DOCKER_BUILDKIT=1 docker build -f experimental.Dockerfile . -t eiladin/go-simple-startpage:$(tag)

classicBuild: ## Build docker eiladin/go-simple-startpage:$(tag)
	docker build -f classic.Dockerfile . -t eiladin/go-simple-startpage:$(tag)

publish: build ## Push docker eiladin/go-simple-startpage:$(tag) to docker hub
  docker push eiladin/go-simple-startpage:$(tag)
  