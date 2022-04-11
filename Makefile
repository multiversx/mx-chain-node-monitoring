SHELL := $(shell which bash)

debugger := $(shell which dlv)

.DEFAULT_GOAL := help

# #########################
# Base commands
# #########################

test:
	@echo "  >  Running unit tests"
	go test -cover -race -coverprofile=coverage.txt -covermode=atomic -v ./...



# #########################
# Manage Node Monitoring locally
# #########################

cmd_dir = cmd/node
binary = node-monitoring

help:
	@echo -e ""
	@echo -e "Make commands:"
	@grep -E '^[a-zA-Z_-]+:.*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":"}; {printf "\t\033[36m%-30s\033[0m\n", $$1}'
	@echo -e ""

build:
	cd ${cmd_dir} && \
		go build -o ${binary} -gcflags='all=-N -l'

api_type="notifier"
run: build
	cd ${cmd_dir} && \
		./${binary}

runb: build
	cd ${cmd_dir} && \
		(./${binary} & echo $$! > ./${binary}.pid)

kill:
	kill $(shell cat ${cmd_dir}/${binary}.pid)

debug: build
	cd ${cmd_dir} && \
		${debugger} exec ./${binary}

debug-ath:
	${debugger} attach $$(cat ${cmd_dir}/${binary}.pid)
