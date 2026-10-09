Вот ещё три варианта, которые **естественно** используют весь ваш стек, но отличаются бизнес-логикой и паттернами. Каждый разбит на реалистичные шаги под 6–12 часов работы.

---
### 🔹 Вариант 1: `URL Shortener + Async Click Analytics`
**Суть:** Классический сокращатель ссылок. Редиректы обслуживаются мгновенно из кеша, клики асинхронно уходят в Kafka, воркер агрегирует статистику и пишет в PostgreSQL.

| Технология | Роль |
|------------|------|
| **REST** | `POST /shorten`, `GET /:code` (редирект), `GET /stats/:code` |
| **Redis** | `SET cache:{code} {original_url} EX 86400`, `HINCRBY stats:{code} clicks 1` |
| **Kafka** | Топик `click.events`. Producer шлёт `{code, ip, ts, ua}` при каждом редиректе |
| **Consumer** | Читает пачками, группирует по `code+date`, делает `INSERT ... ON CONFLICT UPDATE` в PG |
| **PostgreSQL** | `urls(id, code, original_url, created_at)`, `daily_stats(code, date, clicks)` |
| **Логирование** | `slog` с `request_id`, `code`, `cache_hit=true/false`, `duration_ms` |
| **Тесты** | Моки `Store`, `Cache`, `Broker`. Проверяем: генерацию кода, fallback при miss в кеше, валидацию URL |

**📅 План (3 дня)**
1. `docker-compose`, базовый роутер, генератор коротких кодов, `POST /shorten` → сохранение в PG + кеш
2. `GET /:code` → чтение из Redis → если miss → PG → редирект. Producer Kafka на каждый редирект
3. Consumer (batch 5 сек/30 событий) → апдейт PG. Unit-тесты хендлеров, `docker compose up`, README

**⚖️ Честная оценка:** `⭐⭐☆☆☆` (самый дружелюбный для новичка). Логика линейная, нет сложных состояний. Kafka здесь работает как "fire-and-forget" буфер, что прощает ошибки.

---
### 🔹 Вариант 2: `Background Job Tracker`
**Суть:** REST-сервис для отправки задач в фон (например, "сгенерировать отчёт", "экспорт CSV"). Статус отдаётся мгновенно из Redis, реальная обработка идёт асинхронно, результат пишется в PG.

| Технология | Роль |
|------------|------|
| **REST** | `POST /jobs`, `GET /jobs/:id/status`, `GET /jobs` |
| **Redis** | `HSET job:{id} status=pending ttl=3600`, кеш для `GET /status` (TTL 10с) |
| **Kafka** | Топик `jobs.queue`. Producer кладёт `{"id":"...","type":"export"}` |
| **Consumer** | Читает, эмулирует обработку (`time.Sleep` + рандомный статус `success/failed`), пишет итог в PG, обновляет Redis |
| **PostgreSQL** | `jobs(id, type, payload, status, result, created_at, updated_at)` |
| **Логирование** | Уровни `info` (создание), `warn` (повторная отправка), `error` (паника воркера) |
| **Тесты** | Валидация типа задачи, переходы статусов, обработка дублей `POST` |

**📅 План (3 дня)**
1. Схема БД, `POST /jobs` → генерация ID, сохранение в PG, запись в Redis, отправка в Kafka
2. Consumer: статус `pending → processing → done/failed`. Обновление Redis + PG. `GET /status` с fallback
3. Тесты статусной машины, моки брокера/кеша, Docker Compose, документация

**⚖️ Честная оценка:** `⭐⭐⭐☆☆`. Требует аккуратной работы с состояниями. Отлично показывает паттерн "worker queue", но новичку легко запутаться в race conditions при обновлении статусов. Решается через простые `context` и отсутствие конкурентных апдейтов одного джоба.

---
### 🔹 Вариант 3: `Event Router & Alert Logger`
**Суть:** Сервис принимает события, фильтрует их по критичности, дедуплицирует всплески, маршрутизирует в разные топики Kafka, а критические алерты сохраняет в PG для дашборда.

| Технология | Роль |
|------------|------|
| **REST** | `POST /events` (`{type, severity, message, source}`) |
| **Redis** | Дедупликация: `SETEX dedup:{md5} 60 1`. Кеш последних 50 критических событий |
| **Kafka** | Два топика: `events.normal`, `events.critical`. Producer роутит по `severity` |
| **Consumer** | Читает только `events.critical`, пишет в PG `alerts(id, source, message, seen_count, created_at)` |
| **PostgreSQL** | Хранит алерты, используется для `GET /alerts?limit=20` |
| **Логирование** | `slog` с `severity`, `dropped=true/false`, `route=normal|critical` |
| **Тесты** | Логика роутинга, дедуп-окно, валидация payload, моки |

**📅 План (3 дня)**
1. Валидация `POST /events`, хеш для дедупа, Redis check/set, роутинг в нужный Kafka topic
2. Consumer `critical` → batch insert в PG. `GET /alerts` → Redis cache → fallback PG
3. Тесты роутинга и дедупа, Docker Compose, `curl` примеры, README

**⚖️ Честная оценка:** `⭐⭐⭐☆☆`. Самая "инженерная" задача из трёх. Учит фильтрации, дедупликации и разделению потоков данных. Легко расширить до production-паттерна "alert pipeline".

---
### 📊 Сравнение по ключевым метрикам
| Критерий | URL Shortener | Job Tracker | Event Router |
|----------|---------------|-------------|--------------|
| **Сложность логики** | Низкая | Средняя | Средняя |
| **Риск застрять** | Минимальный | Race conditions при статусе | Дедуп-логика + роутинг |
| **Польза для портфолио** | Классика, легко объяснить | Показывает async/worker | Показывает event-driven архитектуру |
| **Идеально для** | Первого знакомства со стеком | Понимания очередей и статусов | Фильтрации, дедупа, маршрутизации |

---
### 🛡️ Универсальные правила, чтобы уложиться в 3 дня
1. **Не пишите свой пул соединений** → используйте `pgxpool` и `go-redis` из коробки.
2. **Не делайте graceful shutdown с drain Kafka** → достаточно `signal.Notify` + `cancel()` контекста.
3. **Не тестируйте Docker/Compose в `go test`** → тестируйте только чистый Go-код через интерфейсы.
4. **Не используйте ORM** → сырые запросы `pgx` или `sqlc` (если уже знакомы). ORM съест 30% времени на дебаг.
5. **Фиксируйте версии** в `go.mod` и `docker-compose.yml` (Kafka `3.7`, Postgres `16-alpine`, Redis `7-alpine`).

---
### 📦 Что дальше?
Если скажете, какой вариант ближе по духу, я могу:
1. Сгенерировать готовый `docker-compose.yml` с Kafka (Kraft mode, без Zookeeper), PG, Redis и healthcheck-ами
2. Дать стартовую структуру папок + 3 интерфейса (`Store`, `Cache`, `Broker`) и примеры моков
3. Показать минимальный `main.go` с инициализацией `slog`, роутером и graceful shutdown

Какой проект берём в работу?