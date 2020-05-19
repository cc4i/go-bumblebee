

all: gate-src air-src spy-src gql-src bee-src


gate-src:
	echo "Testing gate ..."
	@cd gate && make test

air-src:
	echo "Testing air ..."
	@cd air && make test

spy-src:
	echo "Testing spy ..."
	@cd spy && make test

gql-src:
	echo "Testing gql ..."
	@cd gql && make test

bee-src:
	echo "Testing bee ..."
	@cd bee && ng test --watch=false
