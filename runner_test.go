package taskpusher

import (
	"fmt"
	"testing"
	"time"
	"strconv"
)

type SlowTask struct {
	Sleep  time.Duration
	status int
	ID     string
}

func (t *SlowTask) Run(s chan string) {
	t.status = StateRunning
	time.Sleep(t.Sleep)
	t.status = StateSuccessful
	s <- t.UID()
	
}

func (t *SlowTask) String() string {
	return fmt.Sprintf("%d", t.ID)
}

func (t *SlowTask) Status() int {
	return t.status
}

func (t *SlowTask) UID() string {
	return t.ID
}

func TestRunner(t *testing.T) {

	man := NewManager(4)
	go man.Run()
	
	for i := 1; i <= 10; i++ {
		st := &SlowTask{
			ID: "b" + strconv.Itoa(i),
			Sleep: time.Millisecond*200,
		}
		man.Add(st)
	}

	time.Sleep( 180 * time.Millisecond)
	
	if len(man.Completed()) != 4 {
		t.Errorf("Completed list should be %d was %d", 4, len(man.Completed()))
	}
	
	for i := 1; i <= 10; i++ {
		st := &SlowTask{
			ID: strconv.Itoa(i),
			Sleep: time.Millisecond*100,
		}
		man.Add(st)
	}

	//time.Sleep(time.Second*2)
	
	
	man.Close()
	
	if len(man.Completed()) != 20 {
		t.Error("Wrong completed list")
	}
	
	
	fmt.Println("Exit")

}