cdir = $(shell pwd)

up:
	sudo -S docker-compose -f ./cicd/docker-compose.yml up -d --build

down: shutdown clean

shutdown:
	sudo -S docker-compose -f ./cicd/docker-compose.yml down

clean:
	sudo docker rmi $(sudo docker images | grep '<none>' | awk '{print $3}')

.PHONY: up down