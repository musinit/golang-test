.PHONY: migrate-new
migrate-new:
	sql-migrate new --config="dbconfig.yml" --env="default"

.PHONY: migrate-up
migrate-up:
	sql-migrate up --config="dbconfig.yml" --env="default" --dryrun
	sql-migrate up --config="dbconfig.yml" --env="default"

.PHONY: migrate-down
migrate-down:
	sql-migrate down --config="dbconfig.yml" --env="default" --dryrun
	sql-migrate down --config="dbconfig.yml" --env="default"