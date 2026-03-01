COMPILER=go1.24.13
SHORTENER_PATH=./cmd/shortener
SHORTENER_BIN=shortener
INCREMENT=2

help:
	@echo "Доступные команды:"
	@echo "    make clean"
	@echo "        Удаление сгенерированных файлов."
	@echo "    make test"
	@echo "        Тестирование приложения."
	@echo "    make autotest"
	@echo "        Запуск автотестов."
	@echo "        Номер инкремента указывается в переменной INCREMENT"
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

autotest: build
	@echo "Запуск автотестов:"
	./shortenertest -test.v -test.run=^TestIteration$(INCREMENT)$$ -binary-path=$(SHORTENER_PATH)/$(SHORTENER_BIN)

run: clean test
	$(COMPILER) run $(SHORTENER_PATH)/main.go

build: clean test
	$(COMPILER) build -o $(SHORTENER_PATH)/$(SHORTENER_BIN) $(SHORTENER_PATH)/*.go
