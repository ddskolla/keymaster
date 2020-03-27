all: clean test compile compress
push: all push

clean:
	rm -rf build

compile:
	GOARCH=amd64 GOOS=linux go build -o ./build/api-linux-x64 ./cmd/api
	GOARCH=amd64 GOOS=linux go build -o ./build/cli-linux-x64 ./cmd/cli
	GOARCH=amd64 GOOS=darwin go build -o ./build/cli-darwin-x64 ./cmd/cli
	GOARCH=amd64 GOOS=windows go build -o ./build/cli-win-x64.exe ./cmd/cli
	GOARCH=amd64 GOOS=linux go build -o ./build/ip-oracle-sign-linux-x64 ./cmd/ip_oracle/sign
	GOARCH=amd64 GOOS=linux go build -o ./build/ip-oracle-verify-linux-x64 ./cmd/ip_oracle/verify

compress:
	(cd build; zip keymaster-api.zip api-linux-x64)
	(cd build; zip ip-oracle-sign.zip ip-oracle-sign-linux-x64)
	(cd build; zip ip-oracle-verify.zip ip-oracle-verify-linux-x64)

push:
	(cd build; aws s3 cp keymaster-api.zip s3://keymaster-api/keymaster-api-${VERSION}.zip)

test:
	go test -v ./...

