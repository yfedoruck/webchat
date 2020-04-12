dev:
	@docker-compose down && \
		docker-compose \
			-f docker-compose.yml \
			-f docker-compose.dev.yml \
			up -d --remove-orphans --build \
			&& docker-compose logs
hot:
	docker exec webserver go install github.com/yfedoruck/webchat/cmd/chat && \
	docker-compose restart webchat

web:
	@docker stop webserver && \
		docker-compose \
			-f docker-compose.yml \
			-f docker-compose.dev.yml \
			build webchat && \
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