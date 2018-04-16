help:
	@echo ""
	@echo "usage: make COMMAND"
	@echo ""
	@echo "Commands:"

	@echo "  build-dev"
	@echo "  run-dev a=INT w=INT (a - first octet in ip, w - count workers)"
	@echo "  stop-dev"

build-dev:
	@docker build ./ -t micro_headers_dev

run-dev:
	@docker run --rm -d --name micro_headers_dev_con micro_headers_dev app $(a) $(w)
	@docker logs -f micro_headers_dev_con > output/$(a).data 2>output/error.log &

stop-dev:
	@docker stop micro_headers_dev_con
