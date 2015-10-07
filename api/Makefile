build:
	go build -ldflags "-X main.buildstamp=`date -u '+%Y-%m-%dT%I:%M:%S'` -X main.githash=`git rev-parse HEAD`"

run: build
	./api
