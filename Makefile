include .env
export

# Путь текущей локальной директории для контейнера с базой данных 
# $(shell pwd) - выполни shell команду pwd (print working directory) - путь текущей директории
export PROJECT_ROOT=$(shell pwd)

# Поднимаем окружение для проекта 
env_up:
#Путь текущей локальной директории для контейнера с базой данных 
	@docker compose up -d todoapp-postgres

# Остановка окружения
env_down:
	@docker compose down todoapp-postgres

# Перезапуск базы данных(с полным удалением)
env-cleanup:
# получение подтверждения (read -p) песатает в консоль "передаваемое сообщение"
# создается переменная ans в которую будет класться ответ пользователя
# далее запускается условное ветвление на уровне shell команд 
# (rm) удаление на уровне shell, (-rf) рекурсивное удаление (принудительно без дополнительных подтверждений)
# echo == println в golang, fi - обозначает конец условного ветвления на уровне shell
	@read -p "Очистить все volume файлы окружения? Опасность утери данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todoapp-postgres port-forwarder && \
		rm -rf ${PROJECT_ROOT}/out/pgdata \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

#запуск контейнера с port-forward
env_port_forward:
	@docker compose up -d port-forwarder

#остановка контейнера с port-forward
env_port_close:
	@docker compose down port-forwarder

# Создаем миграцию
# команда make migrate_create seq=name_migration
migrate_create:
# в начале проводится проверка, что бы значение переменной было задано
# if [ -z "$(seq)" ] - если значение seq не задано
# run так как не нужно что-бы сервис долго жил (только выполнит и завершится)
# --rm сразу удалит контейнер после остановки
# -seq читаем из консоли
# команда make migrate_create seq=name_migration
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр -seq. Пример: make migrate_create seq=name_migration"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

# Применение миграций внутрь docker контейнера up
migrate_up:
	@make migrate_action action=up

# Применение миграций внутрь docker контейнера down
migrate_down:
	@make migrate_action action=down

# Вернуть последнюю работающую миграцию внутрь docker контейнера
migrate_fix:
	@make migrate_action action="force 1"

#  Передаем операцию с миграцией через динамическую переменную
# todoapp-postgres - название сервиса передается в качестве изолированной сети подключения к нашей базе данных
# так как они работают в общей переменной - compose файл передаст это значение самостоятельно 
# ? - передается как query параметр, sslmode=disable - отключает обязательное шифрование 
migrate_action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр action. Пример: make migrate_action action=operation_migrate"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
	-path /migrations \
	-database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@todoapp-postgres:5432/$(POSTGRES_DB)?sslmode=disable \
	"$(action)"

#операция по запуску приложения
todoapp_run:
	@export LOGGER_FOLDER="${PROJECT_ROOT}/out/logs" && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/todoapp/main.go