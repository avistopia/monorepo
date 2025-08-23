PROJECTS := arithland-bank


.PHONY: deps
deps:
	pip install -r requirements/base.txt

.PHONY: test-deps
test-deps:
	pip install -r requirements/test.txt

.PHONY: dev-deps
dev-deps:
	pip install -r requirements/dev.txt

.PHONY: install-lint
install-lint: dev-deps
	cd .. && pre-commit install

.PHONY: lint
lint: test-deps
	cd .. && pre-commit run -a

.PHONY: lint-fix
lint-fix: test-deps
	cd .. && pre-commit run -a

.PHONY: test
test: test-deps
	@for project in $(PROJECTS); do \
		echo "Running tests in $$d"; \
		pytest $$project --disable-pytest-warnings --verbose --reuse-db --nomigrations; \
	done
