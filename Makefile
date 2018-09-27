
export PGPASSWORD ?= $(POSTGRES_PASSWORD)
export DB_ENTRY ?= psql -h $(DB_HOST) -p 5432 -U $(POSTGRES_USER)

# TODO: pass unseal keys with comma seperated
export INIT_KEY ?= 208ec9b501edabd9430580763d3d5707bb7b723ec8aaae85fe43f3501a464ba2

# TODO: pass vault root token
export TOKEN ?= ae89e6aa-3b87-3629-0767-d8de41933347

run:	wait-for-postgres
	# create sub process and 2 log files
	@vault server -config=/vault/config/config.hcl

	## TODO: run unseal for each key
	# @vault operator unseal $(INIT_KEY)

	# @vault login $(TOKEN)


wait-for-postgres:
	while ! nc -zv ${DB_HOST} 5432; do echo waiting for postgresql ..; sleep 5; done;


dev:
	@vault server -config=/vault/config/config.hcl

.PHONY: run dev
