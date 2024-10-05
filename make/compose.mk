compose-up: set-env
	@docker compose -f $(DEVOPS_DIR)/compose/docker-compose.$(ENV).yml --env-file $(ENV_FILE) up -d
compose-down: set-env
	@docker compose -f $(DEVOPS_DIR)/compose/docker-compose.$(ENV).yml --env-file $(ENV_FILE) down

MIGRATION := docker compose -f ${DEVOPS_DIR}/compose/docker-compose.$(ENV).yml --env-file $(ENV_FILE) up migrate
db: set-env
	@docker compose -f $(DEVOPS_DIR)/compose/docker-compose.$(ENV).yml --env-file $(ENV_FILE) up db -d
db-up: set-env
	@export MIGRATION_COMMAND="up" && $(MIGRATION)
db-down: set-env
	@export MIGRATION_COMMAND="down -all" && $(MIGRATION)
db-down1: set-env
	@export MIGRATION_COMMAND="down 1" && $(MIGRATION)
db-clear: set-env
	@echo "Clearing all tables in the database..."
	@docker exec -i ${COMPOSE_PROJECT_NAME}-db \
		psql -q -U "${PG_USER}" -d "${PG_NAME}" -c " \
		DO \$$\$$ \
		DECLARE \
			_schema text; \
			_tables text; \
		BEGIN \
			_schema := current_schema(); \
			SELECT string_agg(quote_ident(tablename), ', ') \
			INTO _tables \
			FROM pg_tables \
			WHERE schemaname = _schema; \
			\
			IF _tables IS NOT NULL THEN \
				EXECUTE 'DROP TABLE IF EXISTS ' || _tables || ' CASCADE'; \
				RAISE NOTICE 'All tables have been dropped.'; \
			ELSE \
				RAISE NOTICE 'No tables found to drop.'; \
			END IF; \
		END \$$\$$;"

.PHONY: compose-up compose-down db db-up db-down db-down1 db-clear
