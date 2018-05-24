# be sure to change this if you want to do a test deploy w/o affecting
# currently deployed version
STACKNAME=upbot

help:
	@cat Makefile | grep "^[a-z]\+:"

all: build package deploy clean output

build:
	go build -ldflags="-w -s"
	zip lambda-up.zip lambda-up

package:
	aws cloudformation package --template-file template.yaml \
		--s3-bucket csps-tie --s3-prefix ${STACKNAME} \
		--output-template-file packaged-template.yaml

deploy:
	aws cloudformation deploy --template-file packaged-template.yaml \
		--stack-name ${STACKNAME} --capabilities CAPABILITY_IAM \
		--role-arn $$(cf.sh cfrole) --tags $$(cat tags.txt | tr '\n' ' ')

output:
	@echo aws cloudformation describe-stacks --stack-name ${STACKNAME} \
		--query 'Stacks[0].Outputs[0].OutputValue' --output text

clean:
	rm -f packaged-template.yaml
	rm -f lambda-up lambda-up.zip

########################################################################
# below here are targets useful for testing

# The token is from my test account and won't work with the current live version
# It will work with the commented out kms entry in the template.
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

nuke:
	aws cloudformation delete-stack --stack-name ${STACKNAME}

sam: build
	sam local start-api

# run locally like in live.
# need to have kms perms and enctyped token set in env var.
run: build
	env kmsEncryptedToken="AQICAHh423AFYImRphKjAn5NbUuA+nfDg025bX68TX19pdEh+AGFcpnT4tYEdRH+Te633nwyAAAAdjB0BgkqhkiG9w0BBwagZzBlAgEAMGAGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMGfkVE+GRIflrMQycAgEQgDPiBUa2Mysl3Q2bWA923O3AnnOf1LPEIqXlFrnIp5rI/KyCcgmCmfULcX3VQbD95W5hlsE=" \
	./lambda-up

# Useful services for proxing from internet to local system.
# https://serveo.net/ or https://ngrok.com/
proxy:
	ssh -R tardus:80:localhost:3000 serveo.net || true
	#ncat -l 3000
