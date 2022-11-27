dev_compose = docker compose --profile dev
build-dev:
	$(dev_compose) build bind-api-dev
dev-up:
	$(dev_compose) down && $(dev_compose) up
dev-down:
	$(dev_compose) down

bench_compose = docker compose --profile bench
bench-up:
	$(bench_compose) down && $(bench_compose) up
bench-down:
	$(bench_compose) down