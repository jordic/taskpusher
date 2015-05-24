package taskpusher

import (
	"fmt"
	"sync"
)

func NewManager(size int) *Manager {

	m := &Manager{
		responses: make(chan string, size+1),
		in:        make(chan Tasker, size),
		quit:      make(chan bool),
		size:      size,
		waiting:   make(map[string]Tasker),
		completed: make(map[string]Tasker),
		running:   make(map[string]Tasker),
	}
	m.init()
	return m
}

// Manager is responsible of managing tasks
type Manager struct {
	responses chan string
	in        chan Tasker
	quit      chan bool
	size      int
	waiting   map[string]Tasker
	running   map[string]Tasker
	completed map[string]Tasker
	wg	sync.WaitGroup
	sync.RWMutex
	
}

// Add a task to the queue
func (m *Manager) Add(t Tasker) {
	m.Lock()
	m.waiting[t.UID()] = t
	m.Unlock()
	m.wg.Add(1)
	m.in <- t
}

func (m *Manager) worker() {
	for j := range m.in {
		//fmt.Println("started", j)
		m.Lock()
		m.running[j.UID()] = j
		delete(m.waiting, j.UID())
		m.Unlock()
		j.Run(m.responses)
	}
}

func (m *Manager) init() {
	// Launch size workers
	for i := 0; i < m.size; i++ {
		go m.worker()
	}

}

// Runs the task manager
func (m *Manager) Run() {

	for {
		select {
		case b, ok := <-m.responses:
			if false == ok {
				return
			}
			fmt.Println("Completed", b)
			m.Lock()
			task := m.running[b]
			m.completed[b] = task
			delete(m.running, b)
			m.Unlock()
			m.wg.Done()
		case <-m.quit:
			fmt.Println("Quitting")
			return
		default:

		}

	}

}

// Closes the task manager, waiting to finish it with
// pending tasks..
func (m *Manager) Close() {
	m.wg.Wait()
	m.quit <- true
	
	
	close(m.in)
	close(m.responses)

}

// Waiting returns tasks waiting to be processed
func (m *Manager) Waiting() map[string]Tasker {
	m.RLock()
	t := m.waiting
	m.RUnlock()
	return t
}

// Completed returns tasks completed
func (m *Manager) Completed() map[string]Tasker {
	m.RLock()
	t := m.completed
	m.RUnlock()
	return t
}

// Running return in process tasks
func (m *Manager) Running() map[string]Tasker {
	m.RLock()
	t := m.running
	m.RUnlock()
	return t
}
