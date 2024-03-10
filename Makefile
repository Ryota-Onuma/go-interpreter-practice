DOCKER_IMAGE_NAME = "go-interpreter-practice"
.PHONY bash:
bash:
	docker build --tag $(DOCKER_IMAGE_NAME) -f Dockerfile .
	docker run -it --rm -v .:/home/app/workspace -w /home/app/workspace $(DOCKER_IMAGE_NAME) bash

.PHONY build:
build:
	docker build --tag $(DOCKER_IMAGE_NAME) -f Dockerfile .
	docker run -it --rm -v .:/home/app/workspace -w /home/app/workspace $(DOCKER_IMAGE_NAME) go build -o go-interpreter-practice -ldflags="-w" .