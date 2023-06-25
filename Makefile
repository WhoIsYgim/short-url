
postgres:
	docker-compose -f docker-compose-postgres.yml up -d

in-memory:
	docker-compose -f docker-compose-in-memory.yml up -d
