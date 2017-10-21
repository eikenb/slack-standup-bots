ECR := $(shell cf.sh ecr)
VERSION := $(shell cat Changelog | head -1 | cut -d' ' -f1)

image:
	docker build -t tie/up .

push: tag
	docker push ${ECR}/tie/up:${VERSION}

deploy:
	docker run --rm -i -e AWS_PROFILE -v ~/.aws:/root/.aws stackmgr \
		deploy tie/up:${VERSION} "$$(cat tags.json)" -s 1 -t resource

login:
	eval $$(aws ecr get-login --no-include-email)

# --

tag:
	docker tag tie/up ${ECR}/tie/up:${VERSION}

repo:
	aws ecr create-repository --repository-name tie/up


