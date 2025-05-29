# 🚦 Параллельное программирование в Go: Atomic Variables

**Тема:** Атомарные переменные и атомарные операции в Go

---

> **“Atomic operations are not magic — they are well-designed CPU instructions for one-variable safety.”**  
> <sub>— Go FAQ: Atomic ops vs Mutex [7]</sub>

---

## 📑 Оглавление

- [1. Введение](#1-введение)
- [2. Теория: атомарные переменные в Go](#2-теория-атомарные-переменные-в-go)
  - [2.1. Что такое атомарность?](#21-что-такое-атомарность)
  - [2.2. Почему атомарные операции важны в многопоточности](#22-почему-атомарные-операции-важны-в-многопоточности)
  - [2.3. Как устроены атомарные операции в Go](#23-как-устроены-атомарные-операции-в-go)
  - [2.4. Пакет sync/atomic: принципы и функции](#24-пакет-syncatomic-принципы-и-функции)
  - [2.5. Основные атомарные операции и их подробное описание](#25-основные-атомарные-операции-и-их-подробное-описание)
    - [2.5.1. Load и Store](#251-load-и-store)
    - [2.5.2. Add](#252-add)
    - [2.5.3. Swap](#253-swap)
    - [2.5.4. CompareAndSwap (CAS)](#254-compareandswap-cas)
    - [2.5.5. atomic.Value](#255-atomicvalue)
  - [2.6. Ограничения и подводные камни атомарных операций](#26-ограничения-и-подводные-камни-атомарных-операций)
  - [2.7. Сравнение с mutex и каналами](#27-сравнение-с-mutex-и-каналами)
- [3. Практические примеры](#3-практические-примеры)
- [4. Описание структуры репозитория](#4-описание-структуры-репозитория)
- [5. Домашние задачи (Concurrency)](#5-домашние-задачи-concurrency)
- [6. Как запускать](#6-как-запускать)
- [7. Best Practices и советы](#7-best-practices-и-советы)
- [8. Список источников](#8-список-источников)

---

## 1. Введение

В современном программировании параллелизм и конкурентность — неотъемлемые части большинства производительных приложений. Однако работа с разделяемыми переменными в многопоточной среде всегда связана с опасностью гонок данных (data race).  
Чтобы избежать некорректных состояний, Go предоставляет разработчику богатый набор инструментов синхронизации: от высокоуровневых каналов до низкоуровневых атомарных операций, реализованных в пакете [`sync/atomic`](https://pkg.go.dev/sync/atomic) [9]. Атомарные операции позволяют реализовывать счетчики, флаги, индикаторы состояния и другие простые структуры без блокировок и затрат на переключение контекста [3][4].

---

## 2. Теория: атомарные переменные в Go

### 2.1. Что такое атомарность?

**Атомарность** — это свойство операции быть неделимой: либо она выполнится полностью, либо не выполнится вовсе.  
В контексте многопоточного программирования атомарная операция гарантирует, что ни одна другая горутина или поток не увидит промежуточного состояния переменной — только старое или новое [7][8].  
Это критически важно: если операция инкремента (например, x++) не атомарна, две горутины могут одновременно “прочитать” x, увеличить и записать — и тогда увеличение произойдет только один раз вместо двух.

**Гонка данных (data race)** — ситуация, когда две или более горутин одновременно обращаются к одной переменной, хотя бы одна из них пишет, и нет должной синхронизации [1][2][5].

---

### 2.2. Почему атомарные операции важны в многопоточности

В многопоточном приложении десятки и сотни горутин могут параллельно читать и записывать одни и те же переменные (например, счетчики, индикаторы готовности, указатели на объекты-конфигурации).  
Если не использовать никакой синхронизации, результат работы таких переменных становится абсолютно непредсказуемым — часть изменений будет потеряна, возможны даже краши или бесконечные циклы [1][2][5].

Традиционно для защиты разделяемых данных применяются **мьютексы** — они “запирают” область кода для эксклюзивного доступа. Но мьютексы — “тяжелый” инструмент: они требуют перехода в режим ядра ОС, могут приводить к блокировкам, задержкам и deadlock [4][6][7].

Атомарные операции, напротив, реализуются на уровне процессора (через специальные инструкции) и работают очень быстро, без блокировок и переключения контекста [3][4][9][10]. Их главный минус — они безопасны только для отдельных простых переменных. Для сложных структур лучше использовать мьютексы [6][7].

---

### 2.3. Как устроены атомарные операции в Go

Пакет [`sync/atomic`](https://pkg.go.dev/sync/atomic) реализует атомарные аналоги основных операций над переменными: чтение (Load), запись (Store), инкремент/декремент (Add), обмен (Swap), условный обмен (CompareAndSwap, CAS), а также работу с интерфейсами (`atomic.Value`) [9].

Под капотом используются аппаратные инструкции, такие как LOCK CMPXCHG на x86, которые делают операцию неделимой для всех потоков и ядер [3][4][9][10].  
Это позволяет реализовать счетчики, флаги, указатели, “горячие” статистики и другие примитивы без блокировок и высокой нагрузки на планировщик Go.

---

### 2.4. Пакет sync/atomic: принципы и функции

Пакет [`sync/atomic`](https://pkg.go.dev/sync/atomic) предоставляет универсальные примитивы для всех основных операций над простыми типами:

- **LoadT(addr *T) T** — атомарно читает значение переменной (T = int32, int64, uint32, uint64, Pointer и др.)
- **StoreT(addr *T, val T)** — атомарно записывает новое значение в переменную
- **AddT(addr *T, delta T) T** — атомарно добавляет delta к переменной (или вычитает, если delta < 0)
- **SwapT(addr *T, new T) (old T)** — атомарно заменяет переменную на новое значение, возвращая старое
- **CompareAndSwapT(addr *T, old, new T) bool** — если текущее значение addr равно old, заменяет на new и возвращает true

**atomic.Value** — отдельный тип для атомарного хранения и подмены значений любого конкретного типа (через интерфейс) [9][20].

---

### 2.5. Основные атомарные операции и их подробное описание

#### 2.5.1. Load и Store

**Load** — атомарно читает значение переменной. Гарантирует, что чтение не будет “размазано” между несколькими записями.

**Store** — атомарно записывает значение в переменную. Гарантирует, что запись будет полностью видна другим горутинам, и никакая другая операция не вмешается в процесс [8][15][16][17][18].

**Пример применения:**
```go
var flag uint32
atomic.StoreUint32(&flag, 1) // Писатель
if atomic.LoadUint32(&flag) == 1 { // Читатель
    // ... действия после установки флага ...
}
```
**Memory barrier:**  
- Store — “release”: все предыдущие действия будут видны до записи.
- Load — “acquire”: все последующие действия не начнутся, пока не завершится чтение [8].

---

#### 2.5.2. Add

**Add** — атомарно увеличивает (или уменьшает) значение переменной.  
Используется для реализации потокобезопасных счетчиков, статистик, индикаторов.

**Пример:**
```go
var counter int64
for i := 0; i < 1000; i++ {
    atomic.AddInt64(&counter, 1)
}
```
**Особенно полезно для:**
- реализация счетчиков без блокировок [11]
- горячие метрики в серверных приложениях

---

#### 2.5.3. Swap

**Swap** — атомарно заменяет значение переменной на новое, возвращая старое.  
Применяется, когда нужно одновременно узнать текущее значение и установить новое без промежуточных состояний.

**Пример:**
```go
var state int32 = 0
prev := atomic.SwapInt32(&state, 1)
fmt.Println("Режим до замены:", prev)
```

---

#### 2.5.4. CompareAndSwap (CAS)

**CompareAndSwap (CAS)** — атомарно сравнивает значение переменной с ожидаемым и, если совпадает, заменяет его на новое. Возвращает true, если замена произошла [6][13].

**Применяется для:**
- реализации lock-free структур данных (стеков, очередей)
- неблокирующих обновлений состояния

**Пример (lock-free стек):**
```go
type Node struct {
    value interface{}
    next  *Node
}
var head unsafe.Pointer // *Node

func push(v interface{}) {
    newNode := &Node{value: v}
    for {
        oldHead := (*Node)(atomic.LoadPointer(&head))
        newNode.next = oldHead
        if atomic.CompareAndSwapPointer(&head, unsafe.Pointer(oldHead), unsafe.Pointer(newNode)) {
            break
        }
    }
}
```

---

#### 2.5.5. atomic.Value

**atomic.Value** — специальный тип для атомарного хранения и замены значений любого конкретного типа (через интерфейс).  
Поддерживает методы Store, Load, Swap, CompareAndSwap [9][20].

**Особенности:**
- Для одного экземпляра Value все значения должны быть одного типа (иначе panic).
- После первого Store нельзя копировать Value (иначе потеряется внутреннее состояние) [20].

**Пример (горячее обновление конфигурации):**
```go
var config atomic.Value
config.Store(initialConfig)

go func() {
    for {
        newCfg := loadNewConfig()
        config.Store(newCfg)
        time.Sleep(time.Minute)
    }
}()

func worker() {
    for {
        cfg := config.Load().(Config)
        handle(cfg)
    }
}
```

---

### 2.6. Ограничения и подводные камни атомарных операций

- **Только одна переменная:**  
  Атомарные операции защищают только одну переменную — нельзя атомарно изменить сразу две переменных без дополнительных блокировок [6].
- **Ограниченный набор типов:**  
  Поддерживаются только базовые типы (int32, int64, uint32, uint64, Pointer и др.).
- **Усложнение кода:**  
  Код с атомарными операциями сложнее для понимания и поддержки, чем с мьютексами — легко ошибиться с порядком или забыть про “каскадную” атомарность [6][20].
- **Не для сложных структур:**  
  Для потокобезопасных структур с несколькими полями лучше использовать мьютекс или высокоуровневые каналы [7][9].

---

### 2.7. Сравнение с mutex и каналами

| Механизм           | Преимущества                            | Недостатки                                   |
|--------------------|-----------------------------------------|----------------------------------------------|
| **Атомарные операции** | Минимальные накладные, нет блокировок | Только одна переменная, сложно для сложных структур [6][7][9] |
| **Мьютекс (Mutex)**    | Легко защищает любую секцию           | Возможны блокировки, deadlock                |
| **Каналы (chan)**      | Высокоуровневая синхронизация         | Более тяжелые по сравнению с atomic          |

**Подробнее:** [6][7][9]

---

## 3. Практические примеры

### 3.1. Пример гонки данных и исправление с atomic

**Проблема гонки:**
```go
var ops uint64
var wg sync.WaitGroup
wg.Add(2)
go func() { for i := 0; i < 1000; i++ { ops++ } wg.Done() }()
go func() { for i := 0; i < 1000; i++ { ops++ } wg.Done() }()
wg.Wait()
fmt.Println("ops:", ops) // < 2000, гонка!
```
**Решение с atomic:**
```go
var ops uint64
var wg sync.WaitGroup
wg.Add(2)
go func() { for i := 0; i < 1000; i++ { atomic.AddUint64(&ops, 1) } wg.Done() }()
go func() { for i := 0; i < 1000; i++ { atomic.AddUint64(&ops, 1) } wg.Done() }()
wg.Wait()
fmt.Println("ops:", ops) // всегда 2000
```
_Подробнее: [21]_

---

### 3.2. Атомарный флаг (индикатор готовности)

```go
var done atomic.Uint32 // с Go 1.19+
go func() {
    // ... работа ...
    done.Store(1)
}()
for done.Load() == 0 {
    // ждем завершения
}
fmt.Println("Готово")
```
_Подробнее: [9]_

---

### 3.3. atomic.Value для обновления конфигурации

```go
var config atomic.Value
config.Store(initialConfig)
go serveRequests(&config)

func reloadConfig(newCfg Config) {
    config.Store(newCfg)
}

func serveRequests(cfg *atomic.Value) {
    for {
        current := cfg.Load().(Config)
        // работать с current
    }
}
```
_Подробнее: [9][20]_

---

## 4. Описание структуры репозитория

```
my_lab/
├── cmd/
│   └── demonstration/
│       └── main.go           # Демонстрация работы счетчиков
├── pkg/
│   ├── concurrency/
│   │   ├── counter.go        # Интерфейс Counter, реализации и тесты
│   │   └── counter_test.go
│   ├── examples/
│   │   ├── counter/
│   │   │   └── main.go       # Пример: счетчик
│   │   └── queue/
│   │       └── main.go       # Пример: очередь
│   ├── internal/
│   │   ├── atomic/
│   │   │   ├── problem.go
│   │   │   ├── solution.go
│   │   │   └── solution_test.go
│   │   ├── barrier/
│   │   │   ├── problem.go
│   │   │   ├── solution.go
│   │   │   └── solution_test.go
│   │   ├── cas/
│   │   │   ├── problem.go
│   │   │   ├── solution.go
│   │   │   └── solution_test.go
│   │   ├── deadlock/
│   │   │   ├── problem.go
│   │   │   ├── solution.go
│   │   │   └── solution_test.go
│   │   ├── locker/
│   │   │   ├── problem.go
│   │   │   ├── solution.go
│   │   │   └── solution_test.go
│   │   ├── mutex/
│   │   │   ├── problem.go
│   │   │   ├── solution.go
│   │   │   └── solution_test.go
│   │   └── queue/
│   │       ├── problem.go
│   │       ├── solution.go
│   │       └── solution_test.go
│   └── tasks/
│       ├── tasks.md
│       ├── task1/
│       │   ├── task1.go
│       │   ├── task1_solution.go
│       │   └── task1_test.go
│       └── task2/
│           ├── task2.go
│           ├── task2_solution.go
│           └── task2_test.go
├── go.mod
└── README.md
```

---

## 5. Дополнительные задачи (Concurrency)

- **Задача 1:** Потокобезопасный UniqueUsers  
  - Формулировка: `pkg/tasks/task1/task1.go`
  - Решение: `pkg/tasks/task1/task1_solution.go`
  - Тесты: `pkg/tasks/task1/task1_test.go`
- **Задача 2:** Потокобезопасная очередь  
  - Формулировка: `pkg/tasks/task2/task2.go`
  - Решение: `pkg/tasks/task2/task2_solution.go`
  - Тесты: `pkg/tasks/task2/task2_test.go`
- **Сводное описание всех задач:** `pkg/tasks/tasks.md`

---

## 6. Как запускать

- **Запустить демонстрацию:**
    ```bash
    go run cmd/demonstration/main.go
    ```
- **Запустить все тесты:**
    ```bash
    go test ./...
    ```
- **Протестировать отдельный модуль:**
    ```bash
    go test ./pkg/internal/atomic/...
    go test ./pkg/tasks/task1/...
    ```

---

## 7. Best Practices и советы

- Используйте `go test -race` на каждом этапе, чтобы гарантировать отсутствие гонок [7][13].
- Для простых чисел, флагов и указателей используйте atomic, а для сложных структур — mutex [6][7].
- Не комбинируйте несколько атомиков для одного логического объекта: это не даст атомарности всей последовательности действий [6].
- Для "горячих" счетчиков atomic быстрее mutex [10][11].
- Всегда покрывайте конкурентный код стресс-тестами с параллелизмом и задержками [1][2][13].
- Документируйте все места использования атомарных операций в проекте!

---

## 8. Список источников

1. T. Tu, X. Liu, L. Song и Y. Zhang, «Understanding Real-World Concurrency Bugs in Go,» *ASPLOS 2019*, стр. 1–14, [DOI:10.1145/3297858.3304069](https://doi.org/10.1145/3297858.3304069)
2. T. Tu, X. Liu, L. Song и Y. Zhang, «Understanding Real-World Concurrency Bugs in Go,» *ISSRE 2023*, стр. 582–592, [DOI:10.1109/ISSRE62328.2024.00061](https://doi.org/10.1109/ISSRE62328.2024.00061)
3. «Атомики в Go: особенности внутренней реализации». Хабр, 2023. [https://habr.com/ru/articles/744822/](https://habr.com/ru/articles/744822/)
4. «Go: жарим общие данные. Атомно, быстро и без мьютексов». Хабр, 2024. [https://habr.com/ru/company/ruvds/blog/840748/](https://habr.com/ru/company/ruvds/blog/840748/)
5. «Погружение в параллелизм в Go». Хабр, 2024. [https://habr.com/ru/articles/840750/](https://habr.com/ru/articles/840750/)
6. «Композиция атомиков в Go». AntonZ.ru, 2024. [https://antonz.ru/atomics-composition/](https://antonz.ru/atomics-composition/)
7. «Go FAQ: Какие операции атомарные? Как насчет мьютексов?» Golang блог, 2019. [https://golang-blog.blogspot.com/2019/02/go-faq-atomic-ops-mutex.html](https://golang-blog.blogspot.com/2019/02/go-faq-atomic-ops-mutex.html)
8. Go Team. «The Go Memory Model». Официальная документация Go, 2022. [https://go.dev/ref/mem](https://go.dev/ref/mem)
9. Go Team. «sync/atomic: атомарные операции в Go». Официальная документация Go, 2025. [https://pkg.go.dev/sync/atomic](https://pkg.go.dev/sync/atomic)
10. Caraveo, R. «The Go 1.19 Atomic Wrappers and why to use them». Medium, 2023. [https://medium.com/@deckarep/the-go-1-19-atomic-wrappers-and-why-to-use-them-ae14c1177ad8](https://medium.com/@deckarep/the-go-1-19-atomic-wrappers-and-why-to-use-them-ae14c1177ad8)
11. Vincent. «Go: How to Reduce Lock Contention with the Atomic Package». A Journey With Go (Medium), 2020. [https://medium.com/a-journey-with-go/go-how-to-reduce-lock-contention-with-the-atomic-package-ba3b2664b549](https://medium.com/a-journey-with-go/go-how-to-reduce-lock-contention-with-the-atomic-package-ba3b2664b549)
12. The Quantum Yogi. «The Curious Case of Go’s Memory Model: Simple Language, Subtle Semantics». Medium, 2025. [https://medium.com/@kanishksinghpujari/the-curious-case-of-gos-memory-model-simple-language-subtle-semantics-4d3f2029988c](https://medium.com/@kanishksinghpujari/the-curious-case-of-gos-memory-model-simple-language-subtle-semantics-4d3f2029988c)
13. Parker, N. «Understanding and Using the sync/atomic Package in Go». Coding Explorations, 2024. [https://www.codingexplorations.com/blog/understanding-and-using-the-syncatomic-package-in-go](https://www.codingexplorations.com/blog/understanding-and-using-the-syncatomic-package-in-go)
14. Parker, N. «Understanding Golang's Atomic Package and Mutexes». Coding Explorations, 2023. [https://www.codingexplorations.com/blog/understanding-golangs-atomic-package-and-mutexes](https://www.codingexplorations.com/blog/understanding-golangs-atomic-package-and-mutexes)
15. Dulitha. «Mastering Synchronization Primitives in Go». HackerNoon, 2023. [https://hackernoon.com/mastering-synchronization-primitives-in-go](https://hackernoon.com/mastering-synchronization-primitives-in-go)
16. Pang. «Is assigning a pointer atomic in Go?» Stack Overflow, 2014. [https://stackoverflow.com/questions/21447463/is-assigning-a-pointer-atomic-in-go](https://stackoverflow.com/questions/21447463/is-assigning-a-pointer-atomic-in-go)
17. Drathier. «Is variable assignment atomic in go?» Stack Overflow, 2016. [https://stackoverflow.com/questions/33715241/variable-assignment-atomic-in-go](https://stackoverflow.com/questions/33715241/variable-assignment-atomic-in-go)
18. api. «Does golang atomic.Load have a acquire semantics?» Stack Overflow, 2019. [https://stackoverflow.com/questions/55909553/does-golang-atomic-load-have-an-acquire-semantics](https://stackoverflow.com/questions/55909553/does-golang-atomic-load-have-an-acquire-semantics)
19. Hugh. «Is there a difference in Go between a counter using atomic operations and
