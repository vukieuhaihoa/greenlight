.PHONY: dev-run
dev-run:
	go run ./cmd/api 


# migrate create -seq -ext .sql -dir ./migrations add_movies_check_constraints
# migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up