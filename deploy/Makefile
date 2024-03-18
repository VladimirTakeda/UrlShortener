ifeq ($(shell test -e '.env' && echo -n yes),yes)
	include .env
endif

include .env
export

# Manually define main variables

ifndef APP_PORT
export APP_PORT = 8000
endif

ifndef APP_HOST
export APP_HOST = 127.0.0.1
endif

# parse additional args for commands

args := $(wordlist 2, 100, $(MAKECMDGOALS))
ifndef args
MESSAGE = "No such command (or you pass two or many targets to). List of possible commands: make help"
else
MESSAGE = "Done"
endif

ROOT_FOLDER = ..

HELP_FUNC = \
	%help; while(<>){push@{$$help{$$2//'options'}},[$$1,$$3] \
	if/^([\w-_]+)\s*:.*\#\#(?:@(\w+))?\s(.*)$$/}; \
    print"$$_:\n", map"  $$_->[0]".(" "x(20-length($$_->[0])))."$$_->[1]\n",\
    @{$$help{$$_}},"\n" for keys %help; \


# Commands
env:  ##@Environment Create .env file with variables
	@$(eval SHELL:=/bin/bash)
	@cp .env.example .env

help: ##@Help Show this help
	@echo -e "Usage: make [target] ...\n"
	@perl -e '$(HELP_FUNC)' $(MAKEFILE_LIST)

db:  ##@Database Create database with docker-compose
	docker-compose -f docker-compose.yml up -d --remove-orphans

run:  ##@Application Run application server
	go run $(ROOT_FOLDER)/cmd --config $(ROOT_FOLDER)/deploy

open_db:  ##@Database Open database console inside docker-image
	docker exec -it shortener_postgres psql -d $(POSTGRES_DB) -U $(POSTGRES_USER)

clean_db: ##@Database Stop database and clean its data
	docker stop shortener_postgres && docker rm shortener_postgres

test:  ##@Testing Test application with gotest (be careful: tests drops all data in database from env)
	make db && go test $(ROOT_FOLDER)/pkg/handler

lint:
	golangci-lint run --issues-exit-code 0 --print-issued-lines=false --out-format code-climate:gl-code-quality-report.json,line-number $(ROOT_FOLDER)/pkg/*

format:
	go fmt $(ROOT_FOLDER)/pkg/*

%::
	echo $(MESSAGE)
