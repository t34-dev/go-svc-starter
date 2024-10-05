SECRET_REPOSITORY := $(CURDIR)/.gitlab-secrets

include-env-file:
	$(if $(wildcard $(ENV_FILE)), \
		$(eval include $(ENV_FILE)), \
		$(warning $(ENV_FILE) does not exist. Skipping.))

include-secret-git-file:
	$(if $(wildcard $(SECRET_REPOSITORY)), \
		$(eval include $(SECRET_REPOSITORY)), \
		$(warning $(SECRET_REPOSITORY) does not exist. Skipping.))

# Set environment
init-env-file:
	$(if $(wildcard .env.$(ENV)), \
		@cp .env.$(ENV) $(ENV_FILE), \
		$(warning .env.$(ENV) does not exist. Skipping copy.))

set-env: init-env-file include-env-file include-secret-git-file
	@$(eval PROJECT_NAME := $(APP_NAME)_$(ENV))
	@$(eval COMPOSE_PROJECT_NAME := $(PROJECT_NAME))


set-gitlab-config: set-env
	@echo "Setting up global Git configuration for GitLab"
	git config --global url."https://oauth2:${GITLAB_TOKEN}@gitlab.com".insteadOf "https://gitlab.com"

set-gitlab-config-local: set-env
	@echo "Setting up local Git configuration for GitLab"
	git config --local url."https://oauth2:${GITLAB_TOKEN}@gitlab.com".insteadOf "https://gitlab.com"
