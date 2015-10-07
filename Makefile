start: check_vars
	rm -rf chevent-web && \
	git clone https://github.com/codeui/chevent-web && \
		cd chevent-web && \
		git clone https://github.com/codeui/chevent && \
		docker build -t updater .

	-@docker ps -a -q | awk '{ print $$1 }' | xargs docker kill > /dev/null
	-@docker ps -a -q | awk '{ print $$1 }' | xargs docker rm > /dev/null

	docker run --name data -v /data/db busybox

	docker run --restart=always -d --name mongo \
		--volumes-from data mongo mongod --smallfiles

	docker run --restart=always -d --name nginx \
		-p 80:80 -v /var/run/docker.sock:/tmp/docker.sock:ro jwilder/nginx-proxy

	docker run --restart=always -d --name updater \
		-e DOMAIN=$(domain) \
		-v /var/run/docker.sock:/var/run/docker.sock:ro \
		updater /go/bin/updater --mongo mongodb://mongo:27017/chevent

all:
	docker build -t front front/
	docker build -t api api/

	-@docker kill front api 2> /dev/null
	-@docker rm front api 2> /dev/null

	docker run --restart=always -d --name front \
		-p 80 -e VIRTUAL_HOST=$(domain),www.$(domain) front

	docker run --restart=always -d --name api \
		-p 80 --link mongo:mongo -e VIRTUAL_HOST=api.$(domain) api /go/bin/api \
		--addr :80 \
		--mongo mongodb://mongo:27017/chevent \
		--hash $(shell git rev-parse HEAD)

domain=$(strip ${DOMAIN})

check_vars:
ifeq ($(domain),)
	$(error DOMAIN env var is not defined")
endif
