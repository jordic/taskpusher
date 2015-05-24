package taskpusher

import (
	"io/ioutil"
	"net/http"
	"time"
	"fmt"
)

// A Task is something that is runnable, and gets and output
type Tasker interface {
	Run(s chan string)
	Status() int
	UID() string
}

const (
	// Task is stopped and awaiting to run
	StateStopped = iota
	// Task is currently in progress
	StateRunning
	// Task has ended succesaful
	StateSuccessful
	// Task has ended with some kind of error
	StateErroneous
)


// A WebTask is an implementation of task, that currently runs
// a remote url
type WebTask struct {
	ID       string
	URL      string
	status   int
	Duration time.Duration
	Error    error
	Response string
	Client	*http.Client
}

// UID is a unique identifier for the task
func (w *WebTask) UID() string {
	return w.ID
}

// Runs the task. Fetches the url.
// The webhandler expressed by url, should reply with a 
// 200 (StatusOK) if not, task is considered erroneous
func (w *WebTask) Run(s chan string) {
	
	if w.Client == nil {
		w.Client = http.DefaultClient
	}


	t := time.Now()
	w.status = StateRunning
	
	resp, err := w.Client.Get(w.URL)
	if err != nil {
		w.Error = err
		w.status = StateErroneous
		w.Duration = time.Now().Sub(t)
		s <- w.UID()
		return
	}
	
	if resp.StatusCode != http.StatusOK {
		w.Error = fmt.Errorf("wrong status response: %s", resp.StatusCode)
		w.status = StateErroneous
		w.Duration = time.Now().Sub(t)
		s <- w.UID()
		return
	}
	
	
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.Error = err
		w.status = StateErroneous
		w.Duration = time.Now().Sub(t)
		s <- w.UID()
		return
	}

	w.status = StateSuccessful
	w.Duration = time.Now().Sub(t)
	w.Response = string(body)
	s <- w.UID()
	return
}


func (w *WebTask) Status() int {
	return w.status
}