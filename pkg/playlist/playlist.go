package playlist

import (
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/slintes/bluesstammtisch/pkg/types"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Process(uri string) (*types.Playlist, error) {

	var csv string
	var err error
	if strings.HasPrefix(uri, "http") {
		csv, err = h.download(uri)
		if err != nil {
			return nil, err
		}
	} else {
		dat, err := ioutil.ReadFile(uri)
		if err != nil {
			return nil, err
		}
		csv = string(dat)
	}
	pl, err := h.convert(csv)
	if err != nil {
		return nil, err
	}
	return pl, nil
}

func (h *Handler) download(uri string) (string, error) {
	if err := h.verifyUrl(uri); err != nil {
		return "", err
	}

	resp, err := http.Get(uri)
	if err != nil {
		return "", fmt.Errorf("error getting playlist csv: %v", err)
	}

	csvWin, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("error reading playlist csv: %v", err)
	}

	// this is windows-1250 encoded...
	dec := charmap.Windows1250.NewDecoder()
	csvUTF, err := dec.Bytes(csvWin)
	if err != nil {
		return "", fmt.Errorf("error decoding playlist csv: %v", err)
	}

	return string(csvUTF), nil
}

func (h *Handler) verifyUrl(uri string) error {
	// make sure the url makes sense, should be in format
	// https://locker.ifttt.com/*/*Handler*.csv?sharing_key=*
	u, err := url.Parse(uri)
	if err != nil {
		return fmt.Errorf("could not parse url: %v", err)
	}

	if !strings.Contains(u.Host, "locker.ifttt.com") {
		return fmt.Errorf("invalid host, no ifttt.com")
	}

	if !strings.Contains(u.Path, "playlist") {
		return fmt.Errorf("invalid path, no playlist")
	}

	if !strings.Contains(u.Path, ".csv") {
		return fmt.Errorf("invalid host, no .csv")
	}

	return nil
}
