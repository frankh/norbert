APP_NAME ?= "norbert"

VERSION ?= `git rev-parse --short HEAD`

IMAGE_NAME = "frankh/${APP_NAME}"

build:
	docker build --rm \
			-t ${IMAGE_NAME}:${VERSION} \
			.

run:
	docker run --rm -p 8000:8000 \
			${IMAGE_NAME}:${VERSION}

publish:
	docker push ${IMAGE_NAME}:${VERSION}

version:
	@echo "${VERSION}"

.PHONY: build publish version
