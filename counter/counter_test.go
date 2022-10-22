package counter

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type suite struct {
	urls     []string
	expected int
}

func CheckoutDummy(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("id")
	switch key {
	case "0": // found 0
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, ``)
	case "1": // found 1
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `gggo Go gggo`)
	case "2": // found 2
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `<script type="language/javascript">alert("Go")</script><pre>Go hello Go</pre>`)
	case "3": // found 3
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `<h1>Go-Go-Go!</h1>`)
	case "4": // found 0 (search is case sensitive)
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `go go go go`)
	case "5": // found 0
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `<h1></h1>`)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func TestWordFinder(t *testing.T) {
	client := &http.Client{}
	c := New(client, 5, 1*time.Second)

	ts := httptest.NewServer(http.HandlerFunc(CheckoutDummy))
	urls := []string{
		ts.URL + "?id=0",
		ts.URL + "?id=1",
		ts.URL + "?id=2",
		ts.URL + "?id=3",
		ts.URL + "?id=4",
		ts.URL + "?id=5",
		ts.URL + "?id=6",
		ts.URL + "?id=7",
		"http://bad-url",
	}

	got := c.Count("Go", urls)
	expected := 6
	if got != expected {
		t.Errorf("Expecting %d, got %d", expected, got)
	}

	ts.Close()
}
