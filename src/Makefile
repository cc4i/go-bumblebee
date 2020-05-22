

all: was-src air-src gql-src front-src


was-src:
	echo "Testing was ..."
	@cd was && make test

air-src:
	echo "Testing air ..."
	@cd air && make test


gql-src:
	echo "Testing gql ..."
	@cd gql && make test

front-src:
	echo "Testing front ..."
	@cd front && yarn test a --watchAll=false
