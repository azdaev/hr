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

func createTasks(taskChan chan Task, tasks int) {
	for i := 0; i < tasks; i++ {
		taskChan <- Task{id: rand.Int(), createdAt: time.Now()}
	}
	close(taskChan)
}

func taskWorker(task *Task) {
	if task.createdAt.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
		task.error = fmt.Errorf("some error occurred, task id: %d; created: %s", task.id, task.createdAt.Format(time.DateTime))
	}
	task.finishedAt = time.Now()
	time.Sleep(time.Millisecond * 150)
}

func main() {
	workers := 5
	tasks := 50
	taskChan := make(chan Task, tasks)
	doneTasksChan := make(chan Task, tasks)
	successTasks := make([]Task, 0)
	errors := make([]error, 0)

	go createTasks(taskChan, tasks)

	wg := new(sync.WaitGroup)
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			// получение тасков
			for task := range taskChan {
				taskWorker(&task)
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
			errors = append(errors, task.error)
			continue
		}
		successTasks = append(successTasks, task)
	}

	fmt.Println("Errors:")
	for _, err := range errors {
		fmt.Println(err)
	}

	fmt.Println("Done tasks:")
	for _, task := range successTasks {
		fmt.Printf("id: %d, created: %s, finished: %s\n",
			task.id,
			task.createdAt.Format(time.DateTime),
			task.finishedAt.Format(time.DateTime))
	}
}
