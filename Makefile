.PHONY: schema
PSQL_USER=mtasa
PSQL_DB=mtasa_hub

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
	migrate -path database/migrations -database "postgres://${PSQL_USER}@localhost:5432/${PSQL_DB}?sslmode=disable" up
