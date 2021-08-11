.PHONY: clean buildImage pushImage version

GO_VER = 1.16.6
BASE_IMAGE_VER = 20.04-base-0216

APP_NAME = webDownload
APP_VERSION = $(shell git describe --tags --always --dirty="-dev")
PKG_PATH = github.com/zcx2001/webDownload
IMAGE_NAME = webDownload
IMAGE_REPO =

buildImage:
	@sed "s/GO_VER/$(GO_VER)/;s/BASE_IMAGE_VER/$(BASE_IMAGE_VER)/" Dockerfile.In > Dockerfile && \
	    go mod vendor && \
		docker build -q \
			-t $(IMAGE_NAME):$(APP_VERSION) \
			--build-arg APP_NAME=$(APP_NAME) \
			--build-arg APP_VERSION=$(APP_VERSION) \
			--build-arg PKG_PATH=$(PKG_PATH) . && \
		docker system prune -f && \
		rm -rf vendor && rm -f Dockerfile

pushImage: buildImage
	@docker tag $(IMAGE_NAME):$(APP_VERSION) $(IMAGE_REPO)/$(IMAGE_NAME):$(APP_VERSION) && \
	    docker tag $(IMAGE_NAME):$(APP_VERSION) $(IMAGE_REPO)/$(IMAGE_NAME):latest && \
		docker push $(IMAGE_REPO)/$(IMAGE_NAME):$(APP_VERSION) && \
		docker push $(IMAGE_REPO)/$(IMAGE_NAME):latest && \
		docker rmi $(IMAGE_REPO)/$(IMAGE_NAME):$(APP_VERSION) && \
        docker rmi $(IMAGE_REPO)/$(IMAGE_NAME):latest

pushBetaImage: buildImage
	@docker tag $(IMAGE_NAME):$(APP_VERSION) $(IMAGE_REPO)/$(IMAGE_NAME):$(APP_VERSION) && \
		docker push $(IMAGE_REPO)/$(IMAGE_NAME):$(APP_VERSION) && \
		docker rmi $(IMAGE_REPO)/$(IMAGE_NAME):$(APP_VERSION)

version:
	@echo $(APP_VERSION)

clean:
	@go clean