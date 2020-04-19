package utils

import (
	"sync"
	"time"
)

type Task func()

type TaskQueue struct {
	allowedQuota     uint32
	currentQuota     uint32
	tasks            []Task
	periodicExecutor *PeriodicFunc
	lock             sync.Locker
}

func CreateTaskQueue(allowedQuota, timeFrame uint32) *TaskQueue {
	instance := &TaskQueue{
		allowedQuota: allowedQuota,
		currentQuota: 0,
		tasks:        make([]Task, 0),
		lock:         &sync.Mutex{},
	}

	instance.periodicExecutor = CreatePeriodic(time.Duration(timeFrame)*time.Second, instance.onQuotaReset)
	instance.periodicExecutor.Start()
	return instance
}

func (queue *TaskQueue) ScheduleTask(task Task) {
	queue.lock.Lock()

	if queue.currentQuota < queue.allowedQuota {
		queue.currentQuota++
		task()
	} else {
		queue.tasks = append(queue.tasks, task)
	}

	queue.lock.Unlock()
}

func (queue *TaskQueue) onQuotaReset() {
	queue.lock.Lock()
	queue.currentQuota = 0

	for queue.currentQuota < queue.allowedQuota && len(queue.tasks) > 0 {
		task := queue.tasks[0]
		queue.tasks = queue.tasks[1:]

		queue.currentQuota++
		task()
	}

	queue.lock.Unlock()
}
