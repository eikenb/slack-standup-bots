ECR := $(shell cf.sh ecr)
VERSION := $(shell head -1 Changelog | cut -d' ' -f1)
TERMV := $(shell grep '^[0-9]\+' Changelog | sed '2q;d' | cut -d' ' -f1)

versions:
	@echo Current: ${VERSION}
	@echo Last: ${TERMV}

image:
	docker build -t tie/up .

push: tag
	docker push ${ECR}/tie/up:${VERSION}

deploy:
	docker run --rm -i -e AWS_PROFILE -v ~/.aws:/root/.aws stackmgr \
		deploy tie/up:${VERSION} "$$(cat tags.json)" -s 1 -t application

term:
	docker run --rm -i -e AWS_PROFILE -v ~/.aws:/root/.aws stackmgr \
		term tie/up:${TERMV} "$$(cat tags.json)" -s 1 -t application

upgrade: term deploy

login:
	eval $$(aws ecr get-login --no-include-email)

# --

tag:
	docker tag tie/up ${ECR}/tie/up:${VERSION}

repo:
	aws ecr create-repository --repository-name tie/up


