cdir = $(shell pwd)

app-up:
	docker-compose -f ./cicd/docker-compose.yml up -d --build ; \
	./cicd/init.sh

app-down:
	docker-compose -f ./cicd/docker-compose.yml down ; \
	docker rmi $$(sudo docker images -a | grep '<none>' | awk '{print $$3}') ; \
	rm -rf /opt/mysql_slave1/* ; \
	rm -rf /opt/mysql_slave2/*

app-reload:
	docker-compose -f ./cicd/docker-compose.yml down ; \
	docker-compose -f ./cicd/docker-compose.yml up -d ; \
    ./cicd/init.sh

prom-up:
	docker-compose -f ./test/monitor/docker-compose.yml up -d --build

prom-down:
	docker-compose -f ./test/monitor/docker-compose.yml down ; \
	docker rmi $$(sudo docker images -a | grep '<none>' | awk '{print $$3}')

up:
	rm -rf /opt/mysql_slave1/* ; \
	rm -rf /opt/mysql_slave2/* ; \
	docker-compose -f ./cicd/docker-compose.yml up -d --build ; \
	./cicd/init.sh
	docker-compose -f ./test/monitor/docker-compose.yml up -d --build ; \

down:
	docker-compose -f ./test/monitor/docker-compose.yml down ; \
	docker-compose -f ./cicd/docker-compose.yml down ; \
	docker rmi $$(sudo docker images -a | grep '<none>' | awk '{print $$3}')


.PHONY: app-up app-down app-reload prom-up prom-down up down