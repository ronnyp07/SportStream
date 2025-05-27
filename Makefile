s.PHONY: help docker-build

DC := docker compose
DCF := -f ./docker-compose.yaml

up:
	${DC} ${DCF} -p pooller up --build --detach

down:
	${DC} ${DCF} -p pooller down --volumes

install-mockgen:
	@if which mockgen ? /dev/null; then \
		echo "mockgen found, skipping installation, "; \
	else \
		echo "mockgen not found, installing..."; \
		go install github.com/golang/mock/mockgen@v1.6.0; \
	fi

view-api-docs:
	@echo "Open http://localhost:8080/swagger/index.html"

api-test-unit:
	@echo "Running unit tests..."
	cd api && go test -v ./... -cover

mocks-api: install-mockgen
	GO111MODULE=on mockgen -source=api/internal/domain/ports/repos/articles.go -destination=api/tests/mocks/repos/articles.go -package=repomocks
