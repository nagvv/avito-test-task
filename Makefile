all:
	echo "run 'make build' or 'make run'"

build:
	docker build --rm -t avito-task .
	
run:
	docker run -it -p 8080:8080 --entrypoint bash avito-task -c ./avito-test-task