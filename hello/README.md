# GO AWS Lambda Hello Function (Serverless Framework)

This directory is a part of a project Cheems-Writes repository created for deploying backend for a portfolio-cum-blog site.

This is demonstration to deploy a simple **Hello** function written in **Go(Golang)** to **AWS Lambda** using the **Serverless framework** on Windows OS(specifically).

## Directory Structure

```
.
├── hello
│   ├── main.go
│   ├── go.mod
│   └── go.sum
├── bin                  # Compiled binaries and zip files (ignored)
├── .serverless          # Serverless deployment artifacts (ignored)
├── Makefile             # Automation script
├── serverless.yml       # Serverless config
└── .gitignore
```

## Prerequisites

- [Go](https://go.dev/doc/install) (Have Golang installed on your Machine)
- [Serverless](https://www.serverless.com/framework/docs/getting-started/) This page will make you install Node.js too, as it is necessary to install serverless.
- Configure your AWS Profile with a payment card and configure aws cli on your local machine too.
- [build-lambda-zip](github.com/aws/aws-lambda-go) Read the Fuc*ing Documentation here and download the zipping tool `build-lambda-zip` because it is compulsory since WINDOWS Doesn't support zip command on itself.

---

## 🔧 serverless.yml (Configure names and destination folders for yourself)

```yaml
service: <your-service-name>
frameworkVersion: '>= 3.38.0'

provider:
  name: aws
  runtime: provided.al2023 # Go provided runtime DON'T CHANGE!
  region: <aws-region>

functions:
  hello:
    handler: bootstrap
    package:
      patterns:
        - bin/*
      artifact: bin/hello.zip
      individually: true
    events:
      - httpApi:
          path: /hello
          method: get
```
---

## 🛠 Makefile (Configure file sources and dest in this one for yourself)

```makefile
build-hello:
	cd hello && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../bin/bootstrap main.go && cd ..
	cd bin && build-lambda-zip -o hello.zip bootstrap && rm bootstrap && cd ..

build: build-hello

deploy: build
	serverless deploy --aws-profile cheems-writes --verbose

clean:
	rm -rf ./bin ./vendor Gopkg.lock bin .serverless

log:
	serverless logs -f hello --stage dev --region ap-south-1 --aws-profile cheems-writes

---

## 🏁 Deployment Steps

```bash
make deploy
```

To fetch logs:

```bash
make log
```

To clean build artifacts:

```bash
make clean
```
---

> This README will help you quickly get up and running with Go + AWS Lambda + Serverless Framework on Windows. All major errors are documented with explanations.