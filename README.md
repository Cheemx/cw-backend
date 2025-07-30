# ⚙️ CW Backend — Golang + AWS Lambda Serverless API

**CW Backend** is a fully serverless backend built using **Golang**, **MongoDB**, and **AWS Lambda**, deployed via the **Serverless Framework**. It powers the Cheems Writes platform with a highly modular, cost-efficient, and scalable architecture.

---

## 🧠 Key Concepts & Learnings

Initially, I encountered several challenges in deploying multiple Lambdas written in Go, including:

- Each AWS Lambda must be an **independent directory** with a `main.go` and `package main`.
- Reusing shared logic across Lambdas required **modular packages** like `shared/` and `models/`.
- Incorrect structuring caused the `go build` to generate invalid archives instead of an executable.

### ✅ Resolved With:

- Each route (e.g., `tbGetAll`) now lives in its own subfolder under `functions/`
- Build scripts set `GOARCH=amd64` and `GOOS=linux` for AWS compatibility
- A custom Makefile handles ZIP packaging and deployment using Serverless

---

## 📁 Example Folder: `functions/tb-crud/tbGetAll/main.go`

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

---

## 📦 Serverless Configuration

```yaml
service: cw-backend
frameworkVersion: '>= 3.38.0'

plugins:
  - serverless-dotenv-plugin

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

---

## 🛠 Makefile for Build Automation

```makefile
build-tb-getAll:
	cd functions/tb-crud/tbGetAll && env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ../../../bin/bootstrap main.go && cd ../../..
	cd bin && file bootstrap && build-lambda-zip -o tbGetAll.zip bootstrap && rm bootstrap && cd ..

build: build-tb-getAll

deploy: build
	serverless deploy --aws-profile cheems-writes --verbose

log:
	serverless logs -f tb-getAll --stage dev --region ap-south-1 --aws-profile cheems-writes

clean:
	rm -rf ./bin .serverless
```

---

## ✅ Final Takeaway

Each AWS Lambda must:
- Be in its **own folder** with `main.go`
- Use **modular shared packages** for logic reuse
- Be built using proper Linux environment flags
- Be deployed using the **Serverless Framework**

This setup enables a true serverless backend for `cheems-writes` with performance, modularity, and cost-efficiency.
