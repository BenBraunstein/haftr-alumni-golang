OUTPUT_LOCAL = main-local
OUTPUT = main
SERVICE_NAME = haftr-alumni-golang
PACKAGED_TEMPLATE = packaged.yaml # will be archived
TEMPLATE = template.yaml
S3_BUCKET := $(S3_BUCKET)
ZIPFILE = lambda.zip

clean:
	rm -f $(OUTPUT_LOCAL)
	rm -f $(OUTPUT)
	rm -f $(ZIPFILE)

.PHONY: install
install:
	go get -t ./...

# compile the code to run in Lambda (local or real)
.PHONY: lambda
lambda:
	GOOS=linux GOARCH=amd64 $(MAKE) main

main:
	go build -o $(OUTPUT) ./cmd/$(SERVICE_NAME)-lambda/main.go

$(ZIPFILE): clean lambda
	zip -9 -r $(ZIPFILE) $(OUTPUT)

.PHONY: build
build: clean lambda

# TODO: Encrypt package in S3 with --kms-key-id
.PHONY: package
package:
	aws s3 cp open-api-integrated.yaml s3://$(S3_BUCKET)/open-api/$(SERVICE_NAME)/open-api-integrated.yaml
	aws cloudformation package --template-file $(TEMPLATE) --s3-bucket $(S3_BUCKET) --output-template-file $(PACKAGED_TEMPLATE)

validate-template:
	cfn-lint -c I -t template.yaml

build-local:
	go build -o $(OUTPUT_LOCAL) ./cmd/$(SERVICE_NAME)/main.go

run: build-local
	@echo ">> Running application ..."
	PORT=8416 \
	MONGO_URI="mongodb+srv://dbUser:MH646HuSDmjHjfqG@cluster0.w4yor.mongodb.net/<dbName>?retryWrites=true&w=majority" \
	DB_NAME=haftr \
	S3_BUCKET=haftr-alumni-photos-dev \
	JWT_SECRET=6a52a3a680d3677883c4a787872ea95b82a5714ae0619472294d834e032e6db327e64867e376dc4a2e6369b6dd1121bb29c05c478ec12190881a8ec1f90c519a \
	./$(OUTPUT_LOCAL)