.PHONY: deps clean build

deps:
	go get -u ./...

clean: 
	rm -rf ./auth-handler
	
build:
	GOOS=linux GOARCH=amd64 go build -o auth-handler/auth-handler ./auth-handler

	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/handlers/token handlers/token.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/handlers/refreshToken handlers/refresh_token.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/handlers/welcome handlers/welcome.go


clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose