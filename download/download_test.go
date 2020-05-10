package download

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	httpserver *httptest.Server
	mux        *http.ServeMux
	mgr        *Manager
)

func download_setup() {
	l, _ := zap.NewDevelopment()
	logger = l.Sugar()
	mux = http.NewServeMux()
	httpserver = httptest.NewServer(mux)
	bus := make(chan struct{}, 3)
	mgr = &Manager{
		Bus: bus,
		Ua:  "XKCLI/test-suite",
	}
	// Override urlfmt to call a local server
	urlfmt = fmt.Sprintf("%s/%%sinfo.0.json", httpserver.URL)
}

func download_teardown() {
	logger.Sync()
	httpserver.Close()
	mgr.Close()
}

func handle(route string, response string, err error) {
	mux.HandleFunc(
		route,
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if err == nil {
				rw.Write([]byte(response))
			} else {
				http.Error(rw, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}),
	)
}

func TestGet(t *testing.T) {
	download_setup()
	defer download_teardown()
	response := `{"month": "1", "num": 1, "link": "", "year": "2006", "news": "", "safe_title": "Barrel - Part 1", "transcript": "", 
	"alt": "Don't we all.", "img": "https://imgs.xkcd.com/comics/barrel_cropped_(1).jpg", "title": "Barrel - Part 1", "day": "1"}`
	handle("/1/info.0.json", response, nil)
	w := mgr.Get(1)
	if w == nil {
		t.Error("Could not download item 1")
	}
	assert.Equal(t, w.ID, 1, "The retreived ID isn't 1")
	assert.Equal(t, len(mgr.Bus), 0, "The channel was not freed")
	assert.Equal(t, w.Transcript, "", "The transcript is not empty")
}

// Test what happens if the server returns a 404 not found html page (yuck!)
func TestGetNotFound(t *testing.T) {
	download_setup()
	defer download_teardown()
	// We don't call handle(), so there is no new route defined.
	w := mgr.Get(3)
	if w != nil {
		t.Errorf("Fond a non-nil result from a 404 response: %v", w)
	}
}

// In case of faulty response from the server, we don't get an object
func TestGetServerError(t *testing.T) {
	download_setup()
	defer download_teardown()
	handle("/666/info.0.json", "", fmt.Errorf("internal"))
	w := mgr.Get(666)
	if w != nil {
		t.Errorf("Fond a non-nil result from a 404 response: %v", w)
	}
}

func TestGetLatestId(t *testing.T) {
	download_setup()
	defer download_teardown()
	response := `{"month": "1", "num": 1337, "link": "", "day": "19", "year": "2038", "news": "", "safe_title": "End of times", "transcript": "",
	"alt": "", "img": "https://imgs.xkcd.com/comics/barrel_cropped_(1).jpg", "title": "y38kes"}`
	handle("/info.0.json", response, nil)
	id := mgr.GetLatestID()
	assert.Equal(t, id, 1337, "The latest ID doesn't correspond to the server response")
}
