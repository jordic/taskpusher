package taskpusher

import (
	"fmt"
	"testing"
	"time"
)

type SlowTask struct {
	Sleep  time.Duration
	status int
	ID     int
}

func (t *SlowTask) Run(s chan int) {
	t.status = StateRunning
	
	time.Sleep(t.Sleep)
	t.status = StateSuccessful
	s <- t.ID
	
}

func (t *SlowTask) String() string {
	return fmt.Sprintf("%d", t.ID)
}

func (t *SlowTask) Status() int {
	return t.status
}

func TestRunner(t *testing.T) {

	man := NewManager(3)

	

	for i := 1; i <= 10; i++ {
		st := &SlowTask{
			ID: i,
			Sleep: time.Millisecond*500,
		}
		man.Add(st)
	}

	go man.Run()
	
	time.Sleep(time.Second*2)
	
	man.Close()
	
	
	fmt.Println("Exit")

}