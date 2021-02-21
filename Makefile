cdir = $(shell pwd)

app-up:
	docker-compose -f ./cicd/dc_app.yml up -d --build

app-down:
	docker-compose -f ./cicd/dc_app.yml down

app-reload: app-down app-up

db-up:
	rm -rf /opt/mysql_master/* ; \
	rm -rf /opt/mysql_slave1/* ; \
	rm -rf /opt/mysql_slave2/* ; \
	docker-compose -f ./cicd/dc_db.yml up -d --build ; \
	./cicd/init.sh

db-down:
	docker-compose -f ./cicd/dc_db.yml down

client-up:
	docker-compose -f ./cicd/dc_client.yml up -d --build

client-down:
	docker-compose -f ./cicd/dc_client.yml down

prom-up:
	docker-compose -f ./test/monitor/docker-compose.yml up -d --build

prom-down:
	docker-compose -f ./test/monitor/docker-compose.yml down

up: db-up app-up prom-up

down: prom-down app-down db-down

clean:
	docker rmi $$(sudo docker images -a | grep '<none>' | awk '{print $$3}')

.PHONY: app-up app-down app-reload prom-up prom-down up down