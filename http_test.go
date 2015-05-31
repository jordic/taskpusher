package taskpusher

import (
	"net/http"
	"testing"
)


func TestApiPath(t *testing.T) {

	a := &Api{
		Prefix: "/asdf",
	}

	r, err := http.NewRequest("GET", "/asdf/api", nil)
	if err != nil {
		t.Error(err)
	}

	res := a.Path(r)
	if res != "/api" {
		t.Errorf("Path should be %s, provided %s", "/api", res)
	}


	a.Prefix = ""
	res = a.Path(r)
	if res != "/asdf/api" {
		t.Errorf("Path should be %s, provided %s", "/asdf/api", res)
	}
	
	a.Prefix = "asdf"
	res = a.Path(r)
	if res != "/asdf/api" {
		t.Errorf("Path should be %s, provided %s", "/asdf/api", res)
	}
	a.Prefix = "/asdf/"
	res = a.Path(r)
	if res != "/api" {
		t.Errorf("Path should be %s, provided %s", "/api", res)
	}
}