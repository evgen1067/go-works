package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type ErrorCount struct {
	mutex sync.Mutex
	value int
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup
	var errorCount ErrorCount

	channel := make(chan Task, len(tasks)) // буферизованный канал с емкостью колва задач
	for _, val := range tasks {
		channel <- val // запись в канал
	}
	close(channel) // закрыли канал

	wg.Add(n) // инкрементим на n

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done() // декрементим на выходе из функции
			for t := range channel {
				err := t()
				if errorCount.value >= m {
					return
				}
				if err != nil {
					errorCount.value++
				}
			}
		}()
	}
	wg.Wait() // блокируем выполнение пока счетчик не будет равен нулю

	if errorCount.value >= m {
		return ErrLimitExceeded
	}

	return nil
}
