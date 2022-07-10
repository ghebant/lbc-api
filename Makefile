dev:
	@echo "Running dev"
	@docker-compose -f docker-compose.dev.yml up 	\
			--force-recreate 						\
			--quiet-pull							\
			--no-color								\
			--remove-orphans 						\
			--timeout 20
	@docker-compose -f docker-compose.dev.yml rm -f

test:
	@echo "Running tests"
	@docker-compose -f docker-compose.test.yml up 	\
			--abort-on-container-exit				\
			--force-recreate 						\
			--no-color								\
			--remove-orphans 						\
			--timeout 20
	@docker-compose -f docker-compose.test.yml down
	@docker-compose -f docker-compose.test.yml rm -f

test-coverage:
	mkdir -p reports
	make test
	go tool cover -func ./reports/coverage.out

.PHONY: all test dev test-coverage