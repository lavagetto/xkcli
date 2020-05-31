package download

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

var urlfmt = "https://xkcd.com/%sinfo.0.json"
var logger *zap.SugaredLogger

// SetLogger sets the logger for the package
func SetLogger(l *zap.SugaredLogger) {
	logger = l
}

// Manager is the container for sending multiple client requests.
type Manager struct {
	Bus chan struct{}
	Ua  string
}

// performs the request to the client
func (d *Manager) request(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", d.Ua)
	return http.DefaultClient.Do(req)
}

// Get fetches data about one comic strip
func (d *Manager) Get(Id int) *WireXKCD {
	// occupy a slot in the channel
	d.Bus <- struct{}{}
	logger.Debug("Started downloading ", Id)
	// Free the slot once execution is done.
	defer func() { <-d.Bus }()
	idURL := ""
	if Id > 0 {
		idURL = fmt.Sprintf("%d/", Id)
	}
	fullURL := fmt.Sprintf(urlfmt, idURL)
	resp, err := d.request(fullURL)
	if err != nil {
		logger.Errorw("Error downloading", "id", Id, "error", err.Error())
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode > 399 {
		logger.Errorw("Error response from server", "strip", fullURL, "code", resp.Status)
		return nil
	}
	wire, err := NewFromWire(resp.Body)
	if err != nil {
		return nil
	}
	logger.Debug("Done downloading ", Id)
	return wire
}

// GetLatestID gets the ID number of the latest XKCD comic strip published.
func (d Manager) GetLatestID() int {
	w := d.Get(0)
	if w == nil {
		logger.Fatalf("Could not get the latest XKCD, bailing out.")
	}
	return w.ID
}

// Close all dangling resources
func (d *Manager) Close() {
	close(d.Bus)
}
