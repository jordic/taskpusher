package taskpusher

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestWebTaskRuns(t *testing.T) {

	c := make(chan int)
	task := WebTask{
		ID:  1,
		URL: "http://www.google.com",
	}

	go task.Run(c)

	i := <-c

	if i != task.ID {
		t.Errorf("Task id, and channel response should be equal")
	}

	if task.Status() != StateSuccessful {
		t.Errorf("Should provide status %d, provided %d", StateSuccessful,
			task.Status())
	}


	// test bad url
	task = WebTask{
		ID:  2,
		URL: "ht://www.google.com",
	}
	go task.Run(c)
	i = <-c
	if task.Status() != StateErroneous {
		t.Errorf("Should provide status %d, provided %d", StateErroneous,
			task.Status())
	}
	

	// test url not found
	task = WebTask{
		ID:  2,
		URL: "http://www.google.com/asdf",
	}

	go task.Run(c)
	i = <-c
	if task.Status() != StateErroneous {
		t.Errorf("Should provide status %d, provided %d", StateErroneous,
			task.Status())
	}

	if !strings.Contains(task.Error.Error(), "404") {
		t.Errorf("404 not found.. %s", task.Error)
	}

}

func TestTaskWebTimeout(t *testing.T) {
	sawSlow := make(chan bool, 1)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello"))
			w.(http.Flusher).Flush()
			sawSlow <- true
			time.Sleep(1 * time.Second)
			return
	}))
	

	//fmt.Print(ts.URL)
	c := make(chan int)
	task := &WebTask{
		ID:  2,
		URL: ts.URL,
		Client: &http.Client{
			Timeout: 300 * time.Millisecond,
		},
	}

	go task.Run(c)

	_ = <-c
	if task.Status() != StateErroneous {
		t.Errorf("Task should return error status, prov: %s", task.Status() )
	}

	fmt.Printf( "%s", task.Error )

	select {
	case <-sawSlow:
		// Ok
	default:
		t.Fatal("Request handling is not calling")
	}
	ts.Close()
}