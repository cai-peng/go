package gopool

import (
	"sync"
)

type GoWorkPool interface {
	RunWork()
	AddTask(task func())
}

type goWorkPool struct {
	threadNum int
	taskQueue chan func()
}

func NewGoWorkPool(threadNum, totalWorks int) GoWorkPool {
	return &goWorkPool{
		threadNum: threadNum,
		taskQueue: make(chan func(), totalWorks),
	}
}

func (g *goWorkPool) RunWork() {
	defer close(g.taskQueue)

	wg := sync.WaitGroup{}
	for i := 0; i < g.threadNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case task := <-g.taskQueue:
					task()
				default:
					return
				}
			}
		}()
	}

	wg.Wait()
}

func (g *goWorkPool) AddTask(task func()) {
	g.taskQueue <- task
}
