.PHONY: build clean deploy

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/create create/*.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/search search/*.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/scrap scrap/*.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

test: 
	go test `go list ./... | grep -v utils` -cover
