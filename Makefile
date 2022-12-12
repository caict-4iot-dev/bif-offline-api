VERSION ?= latest
IMAGE := bif-offline-api:$(VERSION)

build:
	@docker build -t $(IMAGE) .

push:
	@docker push $(IMAGE)

start:
	@docker run -d -p 8888:8888 --name bif-offline-api $(IMAGE)

stop:
	@docker rm -f bif-offline-api

logs:
	@docker logs -f bif-offline-api --tail=200
