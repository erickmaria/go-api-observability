.PHONY: docker-up docker-down traffic-generation

traffic-generation:
	while true; do curl -k https://random-number-api.localhost/rand && echo; sleep 2; done

docker-up:
	@sudo docker compose up --build -d

docker-down:
	@sudo docker compose down --volumes --remove-orphans
