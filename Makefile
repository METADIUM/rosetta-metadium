.PHONY: test build build-local update-tracer update-bootstrap-balances \
    run-mainnet-online run-mainnet-offline run-testnet-online run-testnet-offline
	 

GO_PACKAGES=./services/... ./cmd/... ./configuration/... ./metadium/... 
GO_FOLDERS=$(shell echo ${GO_PACKAGES} | sed -e "s/\.\///g" | sed -e "s/\/\.\.\.//g")
TEST_SCRIPT=go test ${GO_PACKAGES}
PWD=$(shell pwd)
NOFILE=100000

test:
	${TEST_SCRIPT}

build:
	docker build -t rosetta-metadium:latest https://github.com/metadium/rosetta-metadium.git

build-local:
	docker build -t rosetta-metadium:latest .

update-tracer:
	curl https://raw.githubusercontent.com/metadium/go-metadium/master/eth/tracers/internal/tracers/call_tracer.js -o metadium/client/call_tracer.js

update-bootstrap-balances:
	go run main.go utils:generate-bootstrap metadium/genesis_files/mainnet.json rosetta-cli-conf/mainnet/bootstrap_balances.json;
	go run main.go utils:generate-bootstrap metadium/genesis_files/testnet.json rosetta-cli-conf/testnet/bootstrap_balances.json;

run-mainnet-online:
	docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -v "${PWD}/metadium-data:/data" -e "MODE=ONLINE" -e "NETWORK=MAINNET" -e "PORT=8080" -p 8080:8080 -p 30303:30303 rosetta-metadium:latest

run-mainnet-offline:
	docker run -d --rm -e "MODE=OFFLINE" -e "NETWORK=MAINNET" -e "PORT=8081" -p 8081:8081 rosetta-metadium:latest

run-testnet-online:
	docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -v "${PWD}/metadium-data:/data" -e "MODE=ONLINE" -e "NETWORK=TESTNET" -e "PORT=8080" -p 8080:8080 -p 30303:30303 rosetta-metadium:latest

run-testnet-offline:
	docker run -d --rm -e "MODE=OFFLINE" -e "NETWORK=TESTNET" -e "PORT=8081" -p 8081:8081 rosetta-metadium:latest

run-mainnet-remote:
	docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -e "MODE=ONLINE" -e "NETWORK=MAINNET" -e "PORT=8080" -e "GMET=$(gmet)" -p 8080:8080 -p 30303:30303 rosetta-metadium:latest

run-testnet-remote:
	docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -e "MODE=ONLINE" -e "NETWORK=TESTNET" -e "PORT=8080" -e "GMET=$(gmet)" -p 8080:8080 -p 30303:30303 rosetta-metadium:latest

