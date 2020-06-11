### scanner

Пример микросервиса на Go

Требования:
- можно запустить в контейнере
- отсутствие зависимостей от других сервисов
- only data proccessing
- доставка в swarm по коммиту
- маштабируемость

Цель:

- спарсить заголовки по всем адресам ipv4 диапозона. Запрос делать только на 80 порт
- дополнительно сохранять время ответа

Реализация:

- сделать пул клиентов, аргументы - номер первого октета IP адреса (X.\*.\*.*), кол-во воркеров
- Заголовки складывать из стандартного вывода в файл

Пример запуска:

make run a=1 w=10000 > output/1.data
