PARAMS = $(filter-out $@,$(MAKECMDGOALS))

# Ignore `No rule to make target` errors
%:
	@echo ""

# Install development tools
tools:
	@GO111MODULE=off go get -u github.com/codegangsta/gin
	@GO111MODULE=off go get -u github.com/andybar2/team
.PHONY: tools

# Install vendor dependencies
deps:
	@go mod tidy
	@go mod download
	@go mod vendor
.PHONY: deps

# Run development environment
dev:
	@mkdir -p .team/development
	@team env print -s "development" > .team/development/env
	@team files download -s "development" -p ".team/development/firebase-service-account.json"
	@up start
.PHONY: dev

# Deploy production environment
production:
	@read -p "Do you really want to deploy to production? (y/n) " RESP; \
	if [ "$$RESP" = "y" ]; then \
		mkdir -p .team/production; \
		team env print -s "production" > .team/production/env; \
		team files download -s "production" -p ".team/production/firebase-service-account.json"; \
		up deploy production; \
	fi
.PHONY: production

# Stack management
stack:
	@up stack ${PARAMS}
.PHONY: stack
