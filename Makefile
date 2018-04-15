help:
	@echo ""
	@echo "usage: make COMMAND"
	@echo ""
	@echo "Commands:"

	@echo "  build-dev"
	@echo "  run-dev a=INT w=INT"

build-dev:
	@docker build ./ -t micro_headers_dev

run-dev:
	@docker run micro_headers_dev app $(a) $(w)
