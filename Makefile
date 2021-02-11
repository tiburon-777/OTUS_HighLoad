cdir = $(shell pwd)

up:
	docker-compose -f ./cicd/docker-compose.yml up -d --build
	./cicd/init.sh
down: shutdown clean

shutdown:
	docker-compose -f ./cicd/docker-compose.yml down
	sudo docker rmi $$(sudo docker images -a | grep '<none>' | awk '{print $$3}')

clean:
	rm -rf /opt/mysql_master/* ; \
	rm -rf /opt/mysql_slave1/* ; \
	rm -rf /opt/mysql_slave2/* ; \
	sudo docker rmi $$(sudo docker images -a | grep '<none>' | awk '{print $$3}')

.PHONY: up down