# This file has been developed over a variety of projects
# - https://github.com/leafo/streak.club/blob/master/Makefile
# - https://github.com/chillroom/jacr-api
# - https://github.com/teamxiv/growbot-api

.PHONY: schema
PSQL_USER=mta
PSQL_DB=mtahub_dev
MIGRATE=migrate -path database/migrations -database "postgres://${PSQL_USER}@localhost:5432/${PSQL_DB}?sslmode=disable"

reset_schema::
	# kick clients off the database
	psql postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '${PSQL_DB}';"

	# reset schema
	dropdb ${PSQL_DB} --if-exists
	createdb ${PSQL_DB}
	psql ${PSQL_DB} -c 'grant all privileges on database ${PSQL_DB} to ${PSQL_USER};'
	cat schema.sql | psql -U ${PSQL_USER} ${PSQL_DB} > /dev/null
	@echo "Schema has been reset!"

schema.sql::
	pg_dump -s -U ${PSQL_USER} ${PSQL_DB} > schema.sql
	pg_dump -a -t schema_migrations -U ${PSQL_USER} ${PSQL_DB} >> schema.sql
	@echo "Schema has been written to file"

# save a copy of dev database into dev_backup
checkpoint::
	mkdir -p dev_backup
	pg_dump -F c -U ${PSQL_USER} ${PSQL_DB} > dev_backup/$$(date +%F_%H-%M-%S).dump

# restore latest dev backup
restore_checkpoint::
	dropdb ${PSQL_DB}
	createdb ${PSQL_DB}
	psql ${PSQL_DB} -c 'grant all privileges on database ${PSQL_DB} to ${PSQL_USER};'
	pg_restore -U ${PSQL_USER} -d ${PSQL_DB} $$(find dev_backup | grep \.dump | sort | tail -n 1)

migrate::
	${MIGRATE} up
	make schema.sql

migrate_new::
	${MIGRATE} create -ext sql -dir database/migrations ${NAME}