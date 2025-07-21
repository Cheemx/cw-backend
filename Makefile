build-hello:
	cd hello && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../bin/bootstrap main.go && cd ..
	cd bin && build-lambda-zip -o hello.zip bootstrap && rm bootstrap && cd ..

build-login:
	cd login && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../bin/bootstrap main.go && cd ..
	cd bin && build-lambda-zip -o login.zip bootstrap && rm bootstrap && cd ..

build: build-hello build-login

deploy: build
	serverless deploy --aws-profile cheems-writes --verbose

clean:
	rm -rf ./bin ./vendor Gopkg.lock bin .serverless

log:
	serverless logs -f hello --stage dev --region ap-south-1 --aws-profile cheems-writes