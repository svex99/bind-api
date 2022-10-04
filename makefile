run:
	go build -o build/
	sudo build/bind-api

clean:
	sudo rm -f data/bind/conf/named.conf.local
	sudo rm -f data/bind/conf/named.conf.local.bak
	sudo rm -f data/bind/records/*
	sudo rm -f data/bind-api.db
	sudo touch data/bind/conf/named.conf.local

up:
	docker-compose up

down:
	docker-compose down

build-api:
	docker image build -f docks/api.Dockerfile -t bind-api .

run-api:
	docker container run -it --rm -p 2020:2020 -v $(PWD)/data:/go/src/data bind-api

up-test:clean
	mkdir -p handlers/data
	cp -r data/ handlers/

down-test:
	rm -rf handlers/data/

test:
	make up-test
	go test ./... -run $(test) $(pkg)
	make down-test
