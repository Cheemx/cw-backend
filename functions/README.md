# AWS Lambda Deployment with Golang & MongoDB

This directory deploys **6 GET routes as AWS Lambda functions** using **Golang** and **MongoDB** as the backend database.

## 🧠 Issues Faced & Learnings

Initially, I faced an issue while deploying multiple lambdas because:

- In Golang, each **package** (i.e., directory) can only have **one `main` function**.
- Each Lambda requires its **own `main.go`** with a call to `lambda.Start(...)`.
- I was writing multiple route logic under the same folder which caused a conflict during the `go build` process.
- The resulting `bootstrap` file (the executable Lambda expects) was being created as a **regular archive (ar archive)** instead of the expected **ELF executable**.

### ✅ The Fix

After some debugging, I realized:

> Each Lambda function must be in its **own directory** with a **`main` package** and **`main.go` file**.

So now:
- Each Lambda lives in its own subfolder like `tb-crud/tbGetAll/`.
- Shared logic is extracted to common packages like `shared/` and `models/`.
- The final `bootstrap` binary is built properly by setting `GOOS=linux` and `GOARCH=amd64`.

## 📁 Folder Example: `tb-crud/tbGetAll/main.go`

```go
package main

import (
    "context"
    "functions/models"
    "functions/shared"
    "os"

    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)

func getAllTechHandler(c context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
    uri := os.Getenv("MONGO_URI")
    var blogs []models.TechBlog

    return shared.GetAll(uri, "techblogs", &blogs, c, req)
}

func main() {
    lambda.Start(getAllTechHandler)
}
```

## 🛠️ `serverless.yaml` Configuration

```yaml
app: cheems-writes
service: cw-backend
frameworkVersion: '>= 3.38.0'

plugins:
  - serverless-dotenv-plugin

useDotenv: true

custom:
  dotenv:
    path: .env

provider:
  name: aws
  runtime: provided.al2023
  region: ap-south-1
  stage: dev
  environment:
    MONGO_URI: ${env:MONGO_URI}
    JWT_SECRET: ${env:JWT_SECRET}

package: 
  individually: true

functions:
  tb-getAll:
    handler: bootstrap
    package:
      artifact: bin/tbGetAll.zip
      individually: true
    events:
      - httpApi:
          path: /tb/getall
          method: get
```

## 🧾 Old Manual `Makefile`

> My own non-DRY Makefile (don’t judge me — I’ll learn `make` later 😅)

```makefile
build-hello:
	cd hello && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../bin/bootstrap main.go && cd ..
	cd bin && file bootstrap && build-lambda-zip -o hello.zip bootstrap && rm bootstrap && cd ..

build-login:
	cd login && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../bin/bootstrap main.go && cd ..
	cd bin && file bootstrap && build-lambda-zip -o login.zip bootstrap && rm bootstrap && cd ..

build-tb-getAll:
	cd functions/tb-crud/tbGetAll && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../../../bin/bootstrap main.go && cd ../../..
	cd bin && file bootstrap && build-lambda-zip -o tbGetAll.zip bootstrap && rm bootstrap && cd ..

build: build-hello build-login build-tb-getAll build-tb-getOne build-lc-getAll build-lc-getOne build-wc-getAll build-wc-getOne

deploy: build
	serverless deploy --aws-profile cheems-writes --verbose
```

## ✅ Final Takeaway

Each Lambda must:
- Be in its own directory.
- Have its own `main` package and entrypoint.
- Use shared logic from `shared/` or `models/` packages.
- Be built with proper `GOARCH=amd64` and `GOOS=linux` for AWS compatibility.

Happy Learning! 🚀