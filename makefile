# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html

.PHONY: build-up
build-up:
	@echo "Remember to generate config.yaml with make generate-config"
	docker-compose build; docker-compose up -d

.PHONY: generate-config
generate-config: 
	chmod +x scripts/generate_yaml.sh
	./scripts/generate_yaml.sh

.PHONY: help
help:
	@echo "Please use \`make <target>' where <target> is one of"
	@echo "  run-api              to run the api"
	@echo "  run-job-conciliate   to run the conciliate job"
	@echo "  run-job-report       to run the report job"
	@echo "  generate-config      to generate config.yaml"
	@echo "  help                 to show this help message"
