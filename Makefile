all:
	echo "run 'make build' or 'make run'"

build:
	docker build --rm -t avito-task .
	
run:
	docker run -it --entrypoint bash avito-task -c ./avito-test-task