help:
	@echo ""
	@echo "usage: make COMMAND"
	@echo ""
	@echo "Commands:"

	@echo "  build"
	@echo "  enter"
	@echo "  run a=INT w=INT (a - first octet in ip, w - count workers)"
	@echo "  stop"

# компилируем. Для этого собираем контейнер с компилятором Go (если его ещё нет)
# и затем в образ запаковываем бинарник
build:
	@docker build ./ -t pilot114/micro_headers

# зайти в контейнер
enter:
	@docker run -it --rm --name micro_headers_instance pilot114/micro_headers sh

# запускаем, лог ошибок и данные выводим в соотвествующие файлы
run:
	@docker run --rm -d --name micro_headers_instance pilot114/micro_headers /root/app $(a) $(w)
	@docker logs -f micro_headers_instance > output/$(a)_$(w).data 2>output/$(a)_$(w).log &

stop:
	@docker stop micro_headers_instance
