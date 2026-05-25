// ============================================================
// Задача: Планировщик задач с приоритетами  ⚫ Expert
// ============================================================
//
// Вопрос с финальных этапов собеса уровня Staff+.
//
// Реализуй Scheduler — планировщик с:
//   - Приоритетами (High > Medium > Low)
//   - Отменой через context
//   - Зависимостями: задача B стартует только после завершения задачи A
//   - Дедлайнами: задача отменяется если не стартовала до дедлайна
//   - Метриками: время ожидания, время выполнения
//
//   type Scheduler struct { ... }
//
//   func NewScheduler(workers int) *Scheduler
//   func (s *Scheduler) Schedule(task Task) TaskID
//   func (s *Scheduler) Cancel(id TaskID) bool
//   func (s *Scheduler) Wait(id TaskID) error
//   func (s *Scheduler) Shutdown()
//   func (s *Scheduler) Stats() Stats
//
// Проверь:
//   go test -race -v ./...

package main

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type Priority int

const (
	PriorityHigh   Priority = 3
	PriorityMedium Priority = 2
	PriorityLow    Priority = 1
)

type TaskID int64

type Task struct {
	ID        TaskID
	Priority  Priority
	Deadline  time.Time // нулевое = без дедлайна
	DependsOn []TaskID  // ID задач от которых зависим
	Fn        func(ctx context.Context) error
}

type taskState struct {
	task     Task
	ctx      context.Context
	cancel   context.CancelFunc
	done     chan struct{}
	err      error
	queuedAt time.Time
}

type Stats struct {
	Completed int64
	Failed    int64
	Cancelled int64
	Pending   int64
}

type Scheduler struct {
	mu      sync.Mutex
	tasks   map[TaskID]*taskState
	queue   []*taskState // упрощённо — обычный срез, сортируем по приоритету
	workers int
	jobs    chan *taskState
	wg      sync.WaitGroup
	stats   Stats
	nextID  atomic.Int64
	done    chan struct{}
}

// TODO: реализуй NewScheduler
// Подсказка: проинициализируй поля (tasks, jobs, done) и запусти workers горутин s.worker()
// Не забудь про s.wg чтобы Shutdown мог дождаться завершения
func NewScheduler(workers int) *Scheduler {
	s := &Scheduler{
		tasks:   make(map[TaskID]*taskState), // хранит все задачи
		workers: workers,                     // наши горутинки
		jobs:    make(chan *taskState),       // очередь задачь для workers
		done:    make(chan struct{}),         // сигнал завершения работы
	}

	for i := 0; i < workers; i++ { //  запускаем наши. горутины
		s.wg.Add(1)
		go s.worker()
	}
	return s
}

// TODO: реализуй worker — читает из s.jobs, запускает ts.task.Fn(ts.ctx)
// Подсказка: по результату обнови соответствующий счётчик в s.stats (Cancelled если ctx отменён, иначе Failed/Completed)
// Не забудь закрыть ts.done, чтобы Wait разблокировался
// И реагируй на s.done чтобы выйти при Shutdown
func (s *Scheduler) worker() {
	defer s.wg.Done()
	for {
		select {
		case <-s.done:
			return
		case ts := <-s.jobs:
			// запускаем задачи (обработка)
			err := ts.task.Fn(ts.ctx)

			ts.err = err //  сохранили ошибку
			// обновляем статистику
			if ts.ctx.Err() != nil { // контекст отмена
				atomic.AddInt64(&s.stats.Cancelled, 1)
			} else if err != nil { // файловая ошибка
				atomic.AddInt64(&s.stats.Failed, 1)
			} else { // прошло успешно
				atomic.AddInt64(&s.stats.Completed, 1)
			}

			close(ts.done)
		}
	}
}

// TODO: реализуй Schedule
// Подсказка: если есть DependsOn — жди зависимости асинхронно; Deadline — через context
func (s *Scheduler) Schedule(task Task) TaskID {

	id := TaskID(s.nextID.Add(1)) // создаем новый уникальный id задачи

	ctx := context.Background()   // создали контекст
	var cancel context.CancelFunc // объявили но не проиницилизировали

	if !task.Deadline.IsZero() {
		ctx, cancel = context.WithDeadline(ctx, task.Deadline) // если у задачи есть дедлайн , создаем контекст с отменой по времени
	} else {
		ctx, cancel = context.WithCancel(ctx) // отмена контекста вручную  если у задачи нет дедлайна
	}

	task.ID = id // это новый уникальный айди задачи  который сгенерили выше

	ts := &taskState{ // это описание задачи (живое состояние выполнения)
		task:     task,
		ctx:      ctx,
		cancel:   cancel,
		done:     make(chan struct{}),
		queuedAt: time.Now(),
	}

	s.mu.Lock()
	s.tasks[id] = ts // кладем задачу в общую мап с мьютексом потому что она не потоко безопасна
	s.mu.Unlock()

	if len(task.DependsOn) > 0 { //
		go func() {
			for _, depID := range task.DependsOn { // идем по слайсу задач
				if err := s.Wait(depID); err != nil { // ждем пока dependens завершится , или вернет ошибку
					ts.err = err                           // ошибку записываем в нашу  стркутуру состояния
					cancel()                               // отмена контекста
					close(ts.done)                         // закрытие канала
					atomic.AddInt64(&s.stats.Cancelled, 1) // записываем в отменненый контекст
					return
				}
			}

			s.enqueue(ts) // кладем в очередь воркеров
		}()
	} else {
		s.enqueue(ts)
	}

	return id
}

// очередь воркеров 
func (s *Scheduler) enqueue(ts *taskState) {  // принимает состояние задачи 
	s.mu.Lock() // работаем с общим состоянием 
	defer s.mu.Unlock()

	select {
	case <-s.done:
		ts.cancel() // отмена контексата 
		ts.err = context.Canceled // если есть ошибка пишем задача отменена 
		close(ts.done) // закрываем канал 
		atomic.AddInt64(&s.stats.Cancelled, 1) // увиличиваем счетчик 
		return // выходим из  метода 
	default: // если  все гуд идем дальше 
	}

	s.queue = append(s.queue, ts) // добавляем задачи в очередь 

	sort.Slice(s.queue, func(i, j int) bool { // сортирует слайс 
		return s.queue[i].task.Priority > s.queue[j].task.Priority // перестраивает очередь по приоритету 
	})
 
	next := s.queue[0] // первая задача из очереди 
	s.queue = s.queue[1:] // удаляем первый элемент из очереди 

	select {
	case s.jobs <- next:
	case <-s.done:
		next.cancel() // отменяем контекст 
		next.err = context.Canceled  // записываем ошибку 
		close(next.done) // закрываем канал 
		atomic.AddInt64(&s.stats.Cancelled, 1)
	}
}

// TODO: реализуй Cancel
func (s *Scheduler) Cancel(id TaskID) bool {

	s.mu.Lock()
	defer s.mu.Unlock()

	ts, ok := s.tasks[id]
	if !ok {
		return false
	}

	// task уже завершена
	select {
	case <-ts.done:
		return false
	default:
	}

	// отменяем context task
	ts.cancel()

	return true
}

// Wait блокируется до завершения задачи
func (s *Scheduler) Wait(id TaskID) error {
	s.mu.Lock()
	ts, ok := s.tasks[id]
	s.mu.Unlock()
	if !ok {
		return fmt.Errorf("задача %d не найдена", id)
	}
	<-ts.done
	return ts.err
}

// TODO: реализуй Shutdown — останови планировщик и дождись завершения воркеров
// Подсказка: закрытие s.done сигнализирует всем воркерам о выходе
// Не закрывай s.jobs — в него могут писать горутины зависимостей
func (s *Scheduler) Shutdown() {

	// сигнал всем worker goroutine завершиться
	close(s.done)

	// ждём завершения всех workers
	s.wg.Wait()
}

func (s *Scheduler) Stats() Stats {
	return Stats{
		Completed: atomic.LoadInt64(&s.stats.Completed),
		Failed:    atomic.LoadInt64(&s.stats.Failed),
		Cancelled: atomic.LoadInt64(&s.stats.Cancelled),
	}
}

func main() {
	sched := NewScheduler(3)

	// Задача A
	idA := sched.Schedule(Task{
		Priority: PriorityHigh,
		Fn: func(ctx context.Context) error {
			time.Sleep(100 * time.Millisecond)
			fmt.Println("задача A выполнена")
			return nil
		},
	})

	// Задача B зависит от A
	idB := sched.Schedule(Task{
		Priority:  PriorityMedium,
		DependsOn: []TaskID{idA},
		Fn: func(ctx context.Context) error {
			fmt.Println("задача B выполнена (после A)")
			return nil
		},
	})

	// Задача C с дедлайном
	sched.Schedule(Task{
		Priority: PriorityLow,
		Deadline: time.Now().Add(50 * time.Millisecond),
		Fn: func(ctx context.Context) error {
			select {
			case <-time.After(200 * time.Millisecond):
				fmt.Println("задача C выполнена")
				return nil
			case <-ctx.Done():
				fmt.Println("задача C отменена по дедлайну")
				return ctx.Err()
			}
		},
	})

	sched.Wait(idB)
	time.Sleep(100 * time.Millisecond)

	stats := sched.Stats()
	fmt.Printf("Выполнено: %d, Ошибок: %d, Отменено: %d\n",
		stats.Completed, stats.Failed, stats.Cancelled)
}
