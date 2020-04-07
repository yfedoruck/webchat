dev:
	@docker-compose down && \
		docker-compose \
			-f docker-compose.yml \
			-f docker-compose.dev.yml \
			up -d --remove-orphans --build \
			&& docker-compose logs
build:
	docker exec webchat_webchat_1 go build -o /go/bin/webchat github.com/yfedoruck/webchat && \
	docker-compose restart webchat

web:
	@docker stop webserver && \
		docker-compose \
			-f docker-compose.yml \
			-f docker-compose.dev.yml \
			build server && \
		docker start webserver

deb:
	@docker-compose down && \
			docker-compose \
				-f docker-compose.yml \
				-f docker-compose.debug.yml \
				up -d --remove-orphans --build
heroku:
	heroku container:login && \
	heroku container:push --app obscure-fjord-71819 web && \
	heroku container:release --app obscure-fjord-71819 web && \
	heroku logs --tail --app obscure-fjord-71819