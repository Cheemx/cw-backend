LAMBDA_ARCH=amd64
LAMBDA_OS=linux
BIN_DIR=bin

define build_lambda
cd $(1) && env GOARCH=$(LAMBDA_ARCH) GOOS=$(LAMBDA_OS) go build -ldflags="-s -w" ../$(BIN_DIR)/bootstrap $(2).go && cd ..
cd $(BIN_DIR) && build-lambda-zip -o $(3).zip bootstrap && rm bootstrap && cd ..
endef

build-hello:
	$(call build_lambda, hello main hello)

build-login:
	$(call build_lambda, login, main, login)

build-tb-getAll:
	$(call build_lambda, functions/tb-crud, getAll, tbGetAll)

build-tb-getOne:
	$(call build_lambda, functions/tb-crud, getOne, tbGetOne)

build-lc-getAll:
	$(call build_lambda, functions/lc-crud, getAll, lcGetAll)

build-lc-getOne:
	$(call build_lambda, functions/lc-crud, getOne, lcGetOne)

build-wc-getAll:
	$(call build_lambda, functions/wc-crud, getAll, wcGetAll)

build-wc-getOne:
	$(call build_lambda, functions/wc-crud, getOne, wcGetOne)

build: build-hello build-login build-tb-getAll build-tb-getOne build-lc-getAll build-lc-getOne build-wc-getAll build-wc-getOne

deploy: build
	serverless deploy --aws-profile cheems-writes --verbose

clean:
	rm -rf ./bin ./vendor Gopkg.lock bin .serverless

log:
	serverless logs -f hello --stage dev --region ap-south-1 --aws-profile cheems-writes