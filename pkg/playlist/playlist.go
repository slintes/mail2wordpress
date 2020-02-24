package playlist

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/slintes/bluesstammtisch/pkg/types"
	"golang.org/x/text/encoding/charmap"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Process(data string) (*types.Playlist, error) {

	var csv string
	var err error
	if strings.HasPrefix(data, "http") {
		log.Infof("downloading playlist from %s", data)
		csv, err = h.download(data)
		if err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(data, "/") {
		datWin, err := ioutil.ReadFile(data)
		if err != nil {
			return nil, err
		}
		dec := charmap.Windows1250.NewDecoder()
		datUtf, err := dec.Bytes(datWin)
		if err != nil {
			return nil, err
		}
		csv = string(datUtf)
	} else {
		// we already have the playlist as base64 encoded string
		bytes, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return nil, err
		}
		csv = string(bytes)
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
