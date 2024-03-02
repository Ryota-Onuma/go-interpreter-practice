DOCKER_IMAGE_NAME = "lox-by-go"
.PHONY bash:
bash:
	docker build --tag $(DOCKER_IMAGE_NAME) -f Dockerfile .
	docker run -it --rm -v .:/home/app/workspace -w /home/app/workspace $(DOCKER_IMAGE_NAME) bash
