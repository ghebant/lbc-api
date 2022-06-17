dev:
	@echo "Running dev"
	@docker-compose -f docker-compose.dev.yml up 	\
			--build									\
			--force-recreate 						\
			--quiet-pull							\
			--no-color								\
			--remove-orphans 						\
			--timeout 20
	@docker-compose rm -f

test:
	@echo "Running tests"
	@docker-compose -f docker-compose.test.yml up 	\
	    --build \
			--abort-on-container-exit				\
			--force-recreate 						\
			--no-color								\
			--remove-orphans 						\
			--timeout 20
	@docker-compose -f docker-compose.test.yml down
	@docker-compose -f docker-compose.test.yml rm -f

test-coverage:
	mkdir -p reports
	go test -cover ./... -coverprofile ./reports/coverage.out -coverpkg ./...
	go tool cover -func ./reports/coverage.out

all: dev test test-coverage
.PHONY: all test dev test-coverage