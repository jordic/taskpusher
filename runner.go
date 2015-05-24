package taskpusher

import (
	"fmt"
)

func NewManager(size int) *Manager {

	m := &Manager{
		responses: make(chan int, size+1),
		in:   make(chan Tasker, size),
		quit: make(chan bool),
		size:   size,
		amount: 0,
	}
	m.init()
	return m
}

// Manager is responsible of managing tasks
type Manager struct {
	responses chan int
	in        chan Tasker
	quit      chan bool
	size      int
	amount    int
}

func (m *Manager) Add(t Tasker) {
	m.in <- t
}

func (m *Manager) worker() {
	for j := range m.in {
		fmt.Println("started", j)
		j.Run(m.responses)
	}
}

func (m *Manager) init() {
	// Launch size workers
	for i := 0; i < m.size; i++ {
		go m.worker()
	}

}

func (m *Manager) Run() {

	for {
		select {
		case b, ok := <-m.responses:
			if false == ok {
				return
			}
			fmt.Println("Completed", b)
		case <-m.quit:
			fmt.Println("Quitting")
			return
		default:

		}

	}

}

func (m *Manager) Close() {
	m.quit <- true
	close(m.in)
	close(m.responses)

}
