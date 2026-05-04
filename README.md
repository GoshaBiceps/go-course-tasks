# Go Concurrency Tasks

Задачи по конкурентности Go. Каждая задача — отдельный Go-модуль
с подробным описанием, примерами и TODO-стабами.

Уровни сложности: 🟢 Junior · 🟡 Middle · 🔴 Senior · ⚫ Expert

---

## Структура

| # | Раздел | Тема | Задач |
|---|--------|------|-------|
| 01 | [channels](./01-channels/) | Каналы и горутины | 8 |
| 02 | [sync](./02-sync/) | sync.Mutex, Once, Cond, атомики | 7 |
| 03 | [patterns](./03-patterns/) | Паттерны конкурентности | 8 |
| 04 | [classic](./04-classic/) | Классика | 7 |
| 05 | [hard](./05-hard/) | Сложные задачи | 9 |

**Итого: 39 задач**

---

## 01 · Channels

| Задача | Описание | Уровень |
|--------|----------|---------|
| [task01 — Pipeline](./01-channels/task01/) | Числа → квадраты → фильтр чётных через канальный пайплайн | 🟢 |
| [task02 — Fan-Out / Fan-In](./01-channels/task02/) | Распредели задачи по N воркеров, собери результаты | 🟡 |
| [task03 — Done Channel](./01-channels/task03/) | Отмена цепочки горутин через done-канал | 🟡 |
| [task04 — Timeout & Select](./01-channels/task04/) | Запросы к нескольким API, первый ответ выигрывает | 🟡 |
| [task05 — Merge Channels](./01-channels/task05/) | Слить N каналов в один не теряя порядок закрытия | 🟡 |
| [task06 — Bounded Generator](./01-channels/task06/) | Генератор с ограниченным буфером и backpressure | 🟡 |
| [task07 — Ordered Pipeline](./01-channels/task07/) | Параллельная обработка с сохранением исходного порядка | 🟡 |
| [task08 — Tee Channel](./01-channels/task08/) | Раздвоение канала в два получателя без потерь | 🟡 |

## 02 · Sync

| Задача | Описание | Уровень |
|--------|----------|---------|
| [task01 — RWMutex Cache](./02-sync/task01/) | Потокобезопасный TTL-кеш на RWMutex | 🟡 |
| [task02 — Once Singleton](./02-sync/task02/) | Ленивая инициализация соединения с БД через sync.Once | 🟢 |
| [task03 — Semaphore](./02-sync/task03/) | Взвешенный семафор с Acquire/Release | 🟡 |
| [task04 — Barrier](./02-sync/task04/) | Барьер: все горутины ждут пока все дойдут до точки | 🟡 |
| [task05 — Cond: очередь](./02-sync/task05/) | Блокирующая очередь через sync.Cond | 🔴 |
| [task06 — TryLock Mutex](./02-sync/task06/) | Свой мьютекс с TryLock, LockTimeout, LockContext | 🟡 |
| [task07 — Writer-priority RWMutex](./02-sync/task07/) | RWMutex с приоритетом писателей (без starvation) | 🔴 |

## 03 · Patterns

| Задача | Описание | Уровень |
|--------|----------|---------|
| [task01 — Worker Pool](./03-patterns/task01/) | Пул воркеров с Stop и StopNow | 🟡 |
| [task02 — Rate Limiter](./03-patterns/task02/) | Token Bucket: с горутиной и ленивый | 🟡 |
| [task03 — Pub/Sub](./03-patterns/task03/) | Брокер сообщений: подписки, топики, отписка | 🔴 |
| [task04 — Future/Promise](./03-patterns/task04/) | Асинхронное вычисление с ожиданием и цепочкой Then | 🔴 |
| [task05 — Singleflight](./03-patterns/task05/) | Дедупликация параллельных запросов к одному ключу | 🔴 |
| [task06 — Errgroup](./03-patterns/task06/) | Аналог golang.org/x/sync/errgroup с лимитом | 🟡 |
| [task07 — Circuit Breaker](./03-patterns/task07/) | Closed / Open / HalfOpen с метриками | 🔴 |
| [task08 — Debounce & Throttle](./03-patterns/task08/) | Два близких паттерна ограничения вызовов | 🟡 |

## 04 · Classic

| Задача | Описание | Уровень |
|--------|----------|---------|
| [task01 — Dining Philosophers](./04-classic/task01/) | Обедающие философы без дедлока | 🔴 |
| [task02 — Producer-Consumer](./04-classic/task02/) | Производитель-потребитель через каналы и через Cond | 🟡 |
| [task03 — Print In Order](./04-classic/task03/) | Три горутины печатают строго в порядке 1→2→3 | 🟢 |
| [task04 — FooBar](./04-classic/task04/) | Две горутины печатают "FooBar" по очереди | 🟢 |
| [task05 — H2O](./04-classic/task05/) | Горутины-атомы формируют молекулы H2O | 🔴 |
| [task06 — Cigarette Smokers](./04-classic/task06/) | Классическая задача Паттерсона | 🔴 |
| [task07 — Readers-Writers](./04-classic/task07/) | Два варианта: reader-preferring и fair (FIFO) | 🔴 |

## 05 · Hard

| Задача | Описание | Уровень |
|--------|----------|---------|
| [task01 — Sharded Map](./05-hard/task01/) | Потокобезопасная map с шардированием | 🔴 |
| [task02 — Deadlock Puzzles](./05-hard/task02/) | Найди и исправь четыре вида дедлоков и гонок | 🔴 |
| [task03 — Connection Pool](./05-hard/task03/) | Пул соединений с таймаутами и health check | 🔴 |
| [task04 — Scheduler](./05-hard/task04/) | Планировщик с приоритетами, зависимостями, дедлайнами | ⚫ |
| [task05 — Concurrent LRU](./05-hard/task05/) | Потокобезопасный LRU-кеш O(1) | 🔴 |
| [task06 — Actor Model](./05-hard/task06/) | Простая реализация модели акторов (Send / Ask) | ⚫ |
| [task07 — Web Crawler](./05-hard/task07/) | Обход графа с лимитом горутин и дедупликацией | 🔴 |
| [task08 — Retry Job Queue](./05-hard/task08/) | Очередь с экспоненциальным бэкоффом и DLQ | ⚫ |
| [task09 — Parallel ForEach](./05-hard/task09/) | Дженерик-утилита с лимитом и отменой через ctx | 🟡 |

---

## Как запускать

```bash
cd 01-channels/task01
go test -v ./...          # для тестов
go run main.go            # для программ

go test -race ./...       # всегда проверяй на гонки!
```

---

## Советы

- Каждый `main.go` содержит `// TODO:` — что нужно реализовать
- Запускай с `-race` — многие баги видны только детектором гонок
- Ожидаемый вывод указан в заголовке каждого файла
- Тесты покрывают edge cases — сначала прочитай их
