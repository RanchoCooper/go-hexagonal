init:
	@echo "== 👩‍🌾 init =="
	brew install go
	brew install pre-commit
	brew install golangci-lint
	brew upgrade golangci-lint

	@echo "== pre-commit setup =="
	pre-commit install

precommit.rehooks:
	pre-commit autoupdate
	pre-commit install --install-hooks
	pre-commit install --hook-type commit-msg

test:
	@echo "== 🦸‍️ Prepare Dependency =="
	go mod tidy
	@echo "== 🦸‍️ Start Unit Test =="
	go test ./... -v

ci.lint:
	@echo "== 🙆 Start CI Linter=="
	golangci-lint run -v ./... --fix
