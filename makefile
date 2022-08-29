run:
	go build -o build/
	sudo build/bind-api

clean:
	sudo echo -n "" > data/bind/conf/named.conf.local
	sudo rm -f data/bind/conf/named.conf.local.bak
	sudo rm -f data/bind/records/*
	rm -f bind-api.db
