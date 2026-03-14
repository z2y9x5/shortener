COMPILER=
SHORTENER_PATH=./cmd/shortener
SHORTENER_BIN=shortener

help:
	@echo "Доступные команды:"
	@echo "    make clean"
	@echo "        Удаление сгенерированных файлов."
	@echo "    make test"
	@echo "        Тестирование приложения."
	@echo "    make run"
	@echo "        Запуск приложения."
	@echo "    make build"
	@echo "        Сборка приложения."

clean:
	rm -f $(SHORTENER_PATH)/$(SHORTENER_BIN)

test:
	@echo "Форматирование кода:"
	$(COMPILER) fmt ./...
	@echo "Статический анализ:"
	$(COMPILER) vet ./...
	@echo "Тесты приложения:"
	$(COMPILER) test -v -count 1 ./...
	@echo "Покрытие тестами:"
	$(COMPILER) test -count 1 -cover ./...

run: clean test
	$(COMPILER) run $(SHORTENER_PATH)/main.go

build: clean test
	$(COMPILER) build -o $(SHORTENER_PATH)/$(SHORTENER_BIN) $(SHORTENER_PATH)/*.go
