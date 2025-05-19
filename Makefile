#DOCKER_VER=220425
#dynamic so version
DOCKER_VER=220427
DOCKER_NAME=we-mine-digger-assistant

docker:
	DOCKER_BUILDKIT=0 docker build --tag=chenjingdong/$(DOCKER_NAME):$(DOCKER_VER) -f docker/Dockerfile .

docker-no-cache:
	DOCKER_BUILDKIT=0 docker build --tag=chenjingdong/$(DOCKER_NAME):$(DOCKER_VER) --no-cache -f docker/Dockerfile .



docker-dev:
	docker build --tag=chenjingdong/$(DOCKER_NAME)-dev:$(DOCKER_VER) -f docker/Dockerfile.debug .

clean:
	docker-compose -f docker/docker-compose.yml down --rmi all -v --remove-orphans

stop:
	docker-compose -f docker/docker-compose.yml down

save:
	docker save -o ../../docker-tars/$(DOCKER_NAME).tar chenjingdong/$(DOCKER_NAME):$(DOCKER_VER)

sftp:
	#scp ../../docker-tars/$(DOCKER_NAME).tar dell@61.130.65.146:/happy-data/docker/
	scp ../../docker-tars/$(DOCKER_NAME).tar $(USER)@$(HOST):/happy-data/docker/


dev:
	#docker stop we-store-role
	#docker run -d --name we-store-role --net host -p 8888:8888 chenjingdong/role:220425			
	#docker run -d --name we-store-role -p 8888:8888 --add-host=we-store-role-db:10.10.21.253 chenjingdong/role:220425			
	docker-compose -f docker/docker-compose.yml up -d
	
.PHONY: all docker dev stop release clean docker-dev docker sftp everything docker-no-cache

everything: docker save

run-godev:
	#-docker rm go-dev
	#docker run --name go-dev deis/go-dev find / -name gcc
	docker run -it --rm --name go-dev deis/go-dev /bin/bash
	
look-godev:
	docker exec -it go-dev /bin/sh
	
