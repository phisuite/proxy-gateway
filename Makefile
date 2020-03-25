NAME=proxy-gateway

PORT?=80

build:
	docker build -t phisuite/${NAME} .

start:
	docker run --rm -it -p ${PORT}:80 phisuite/${NAME}

publish:
ifdef VERSION
	docker tag phisuite/${NAME} phisuite/${NAME}:${VERSION} && \
	docker push phisuite/${NAME}:latest && \
	docker push phisuite/${NAME}:${VERSION}
else
	echo "VERSION not defined"
endif

debug:
	go run ./src --port ${PORT}
