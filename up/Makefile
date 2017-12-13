ECR := $(shell cf.sh ecr)
VERSION := $(shell head -1 Changelog | cut -d' ' -f1)
PREV:= $(shell grep '^[0-9]\+' Changelog | sed '2q;d' | cut -d' ' -f1)

versions:
	@echo Current: ${VERSION}
	@echo Last: ${PREV}

image:
	docker build -t tie/up .

push: tag
	docker push ${ECR}/tie/up:${VERSION}

deploy:
	cd sceptre && sceptre describe-stack-resources live up > /dev/null \
		&& sceptre update-stack live up \
		|| sceptre create-stack live up

login:
	eval $$(aws ecr get-login --no-include-email)

# --
tag:
	docker tag tie/up ${ECR}/tie/up:${VERSION}

repo:
	aws ecr create-repository --repository-name tie/up

rm-old: login
	aws ecr batch-delete-image --repository-name tie/up \
		--image-ids $(shell \
		aws ecr list-images --repository-name tie/up \
			--filter tagStatus=UNTAGGED --query "imageIds[].imageDigest" \
			--output text | xargs -n 1 | xargs -IXXX printf "imageDigest=XXX ")

ls-old:
	@aws ecr list-images --repository-name tie/up \
		--filter tagStatus=UNTAGGED --query 'imageIds[].imageDigest'
