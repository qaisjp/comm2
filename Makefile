.PHONY: schema

reset_schema::
	dropdb -U mtasa_hub mtasa_hub --if-exists
	createdb -U mtasa_hub mtasa_hub
	cat schema.sql | psql -U mtasa_hub mtasa_hub > /dev/null
	@echo "Schema has been reset!"

schema.sql::
	pg_dump -s -U mtasa_hub mtasa_hub > schema.sql
	@echo "Schema has been written to file"

# save a copy of dev database into dev_backup
checkpoint::
	mkdir -p dev_backup
	pg_dump -F c -U mtasa_hub mtasa_hub > dev_backup/$$(date +%F_%H-%M-%S).dump

# restore latest dev backup
restore_checkpoint::
	dropdb -U mtasa_hub mtasa_hub
	createdb -U mtasa_hub mtasa_hub
	pg_restore -U mtasa_hub -d mtasa_hub $$(find dev_backup | grep \.dump | sort | tail -n 1)

# might not stick with goose
migrate::
	goose --dir=./database/migrations postgres "user=mtasa_hub dbname=mtasa_hub sslmode=disable" up
