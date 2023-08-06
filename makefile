# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html

.PHONY: build-up
build-up:
	@echo "Remember to generate config.yaml with make generate-config"
	docker-compose build; docker-compose up -d

.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down

.PHONY: generate-config
generate-config: 
	chmod +x scripts/generate_yaml.sh
	./scripts/generate_yaml.sh

.PHONY: report-top-cards
report-top-cards:
	docker-compose start reportjob

.PHONY: conciliate-cards
conciliate-cards:
	docker-compose start conciliatejob

.PHONY: help
help:
	@echo "Please use \`make <target>' where <target> is one of"
	@echo "  build-up             to run the api and jobs"
	@echo "  up                   to start the containers in the background"
	@echo "  down                 to stop and remove the containers"
	@echo "  generate-config      to generate the config.yaml file"
	@echo "  report-top-cards     to run the reportJob to generate the top 20 most expensive cards report"
	@echo "  conciliate-cards     to run the conciliateJob to update card prices from Scryfall API"
