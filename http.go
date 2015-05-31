package taskpusher

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Api exposes an http api for handling taskpusher requests.
// Manager should be a manager instance.
// Prefix is used just in case you want to integrate the task manager in 
// something other server, and you want to decide at wich mount point it 
// should be.
type Api struct {
	Manager    Manager
	Prefix string
	Err error
}

// Path returns url, striping prefixed mount point
func (a *Api) Path(r *http.Request) string {
	if a.Prefix != "" {
		if strings.HasSuffix(a.Prefix, "/") {
			a.Prefix = strings.TrimSuffix(a.Prefix, "/")
		}

		return strings.TrimPrefix(r.URL.Path, a.Prefix)
	}
	return r.URL.Path
}

// Serve requests for http Api
func (a *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Strip mountpoint
	url := a.Path(r)

	var err error
	// @todo, chain middlewares..
	switch {
	case url == "/add":
		err = a.Add(w, r)
	case strings.HasPrefix("/task/", url):
		err = a.Status(w, r)
	default:
		err = a.List(w, r)
	}

	if err != nil {
		a.Err = err
		log.Println("error serving request %s", err)
		http.Error(w, "Internal server error", 500)
		return
	}

}

// Add a task to manager http handler
func (a *Api) Add(w http.ResponseWriter, r *http.Request) error {

	cmd := &Command{}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&cmd)
	if err != nil {
		return err
	}
	task := &WebTask{
		ID: NewUID(),
		URL: cmd.Cmd,
	}
	a.Manager.Add( task )
	return nil
}


// Resturn staus for a five task...
func (a *Api) Status(w http.ResponseWriter, r *http.Request) error {
	return nil
}
// List current active tasks in memory
func (a *Api) List(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type Command struct {
	Type  string `json:"type"`
	Cmd   string `json:"cmd"`
	Extra string `json:"extra"`
}

var counter int64 = 1

// Generates a new string uid, to use with a task.
func NewUID() string {
	t := strconv.FormatInt(time.Now().Unix(), 10)
	c := strconv.FormatInt(counter, 10)
	counter++
	return t + "-" + c
}