pub: 
		cd pub && go build ./
		cd pub && ./pub

start-docker:
		sudo docker compose up

stop-docker:
		sudo docker compose down --volumes

sub:
		cd sub/cmd/app && go build ./
		cd sub/cmd/app && ./app


.PHONY: start-docker stop-docker test pub sub