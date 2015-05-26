package taskpusher

import (
	"io/ioutil"
	"time"

	"os"
	"testing"
)

func tempfile() string {
	f, _ := ioutil.TempFile("", "bolt-")
	f.Close()
	os.Remove(f.Name())
	return f.Name()
}

func TestBoltBackend(t *testing.T) {

	tf := tempfile()

	back := &BoltBack{}
	err := back.Open(tf, 0600)
	if err != nil {
		t.Errorf("error opening db: %s", err)
	}

	task := &SlowTask{
		Sleep: time.Millisecond * 500,
		ID:    "1",
	}

	task.SetStatus(StateStopped)
	back.Save(task)

	m := back.Load(StateStopped)
	if m[0].Status() != task.Status() {
		t.Errorf("Unable to load saved tasks")
	}

	task.SetStatus(StateRunning)
	back.Save(task)

	m = back.Load(StateStopped)
	if len(m) != 0 {
		t.Errorf("Loaded tasks should be 0")
	}

	c := back.Load(StateRunning)
	if c[0].Status() != task.Status() {
		t.Errorf("Task status should be %d provided %d", task.Status(),
			c[0].Status())
	}

	os.Remove(tf)

}