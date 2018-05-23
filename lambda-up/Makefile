help:
	@cat Makefile | grep "^[a-z]\+:"

all: build package deploy clean output

build:
	go build -ldflags="-w -s"
	zip lambda-up.zip lambda-up

package:
	aws cloudformation package --template-file template.yaml \
		--s3-bucket csps-tie --s3-prefix upbot \
		--output-template-file packaged-template.yaml

deploy:
	aws cloudformation deploy --template-file packaged-template.yaml \
		--stack-name upbot --capabilities CAPABILITY_IAM \
		--role-arn $$(cf.sh cfrole) --tags $$(cat tags.txt | tr '\n' ' ')

output:
	@echo aws cloudformation describe-stacks --stack-name upbot --query 'Stacks[0].Outputs[0].OutputValue' --output text

test-stand:
	curl --data-urlencode "token=nfTiev5lF9TJ5l3sJ3p6Sqd1" \
		--data-urlencode "text=stand taking it easy" \
		--data-urlencode "channel_name=foo" \
		--data-urlencode "user_name=bar" \
		--data-urlencode "response_url=test" \
		"https://ybyna5t917.execute-api.us-west-2.amazonaws.com/Stage/message"

test-show:
	curl --data-urlencode "token=nfTiev5lF9TJ5l3sJ3p6Sqd1" \
		--data-urlencode "text=show" \
		--data-urlencode "channel_name=foo" \
		--data-urlencode "user_name=bar" \
		--data-urlencode "response_url=test" \
		"https://ybyna5t917.execute-api.us-west-2.amazonaws.com/Prod/message"

clean:
	rm -f packaged-template.yaml
	rm -f lambda-up lambda-up.zip

nuke:
	aws cloudformation delete-stack --stack-name upbot

sam: build
	sam local start-api

run: build
	env kmsEncryptedToken="AQICAHh423AFYImRphKjAn5NbUuA+nfDg025bX68TX19pdEh+AGFcpnT4tYEdRH+Te633nwyAAAAdjB0BgkqhkiG9w0BBwagZzBlAgEAMGAGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMGfkVE+GRIflrMQycAgEQgDPiBUa2Mysl3Q2bWA923O3AnnOf1LPEIqXlFrnIp5rI/KyCcgmCmfULcX3VQbD95W5hlsE=" \
	./lambda-up

proxy:
	ssh -R tardus:80:localhost:3000 serveo.net || true
# https://serveo.net/
#ncat -l 3000
