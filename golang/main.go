package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Task represents a meaninglessness of our life
type Task struct {
	id         int
	createdAt  time.Time // время создания
	finishedAt time.Time // время выполнения
	error      error
}

func CreateTasks(taskChan chan<- Task, nTasks int) {
	for i := 0; i < nTasks; i++ {
		taskChan <- Task{
			id:        rand.Int(),
			createdAt: time.Now(),
		}
	}
	close(taskChan)
}

func (task *Task) Work() {
	if task.createdAt.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
		task.error = fmt.Errorf("some error occurred, task id: %d; created: %s", task.id, task.createdAt.Format(time.DateTime))
	}
	task.finishedAt = time.Now()
	time.Sleep(time.Millisecond * 150)
}

func main() {
	nWorkers := 5
	nTasks := 50
	taskChan := make(chan Task, nWorkers)
	doneTasksChan := make(chan Task)

	go CreateTasks(taskChan, nTasks)

	wg := new(sync.WaitGroup)
	for i := 0; i < nWorkers; i++ {
		wg.Add(1)
		go func() {
			// получение тасков
			for task := range taskChan {
				task.Work()
				doneTasksChan <- task
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(doneTasksChan)
	}()

	for task := range doneTasksChan {
		if task.error != nil {
			fmt.Println(task.error)
			continue
		}
		fmt.Printf("id: %d, created: %s, finished: %s\n",
			task.id,
			task.createdAt.Format(time.DateTime),
			task.finishedAt.Format(time.DateTime),
		)
	}
}
