# Задачи (Concurrency)

## Задача 1

Реализуйте потокобезопасную структуру для подсчёта уникальных пользователей, которые обращаются к сервису параллельно.  
Интерфейс:
```go
type UniqueUsers interface {
    AddUser(id int)
    Count() int
}
```

## Задача 2

Реализуйте потокобезопасную очередь (FIFO) с Push и Pop.  
Интерфейс:
```go
type SafeQueue interface {
    Push(val int)
    Pop() (int, bool)
}
```