all: clean test compile compress

clean:
	rm -rf build

compile:
	GOARCH=amd64 GOOS=linux go build -o ./build/issuing-lambda-linux-x64 ./cmd/issuing_lambda
	GOARCH=amd64 GOOS=linux go build -o ./build/km-linux-x64 ./cmd/km
	GOARCH=amd64 GOOS=darwin go build -o ./build/km-darwin-x64 ./cmd/km
	GOARCH=amd64 GOOS=windows go build -o ./build/km-win-x64.exe ./cmd/km

compress:
	(cd build; zip keymaster-issuing-lambda.zip issuing-lambda-linux-x64)

test:
	go test -v ./...

