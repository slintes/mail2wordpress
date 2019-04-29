package server

import (
	"fmt"
	"github.com/slintes/mail2wordpress/pkg/playlist"
	"github.com/slintes/mail2wordpress/pkg/wordpress"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	plh *playlist.Handler
	wp  *wordpress.Poster
}

func NewServer(plHandler *playlist.Handler, wpPoster *wordpress.Poster) *Server {
	return &Server{
		plh: plHandler,
		wp:  wpPoster,
	}
}

func (s *Server) Serve() {
	// PORT might be injected on Cloud Run
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", s.handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		log.Debug("no POST request, ignoring")
		http.Error(w, "nothing to see here", http.StatusBadRequest)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Errorf("could not read request body: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	playlistUrl := string(b)
	log.Infof("playlist url: %s", playlistUrl)

	pl, err := s.plh.Process(playlistUrl)
	if err != nil {
		log.Errorf("could not process playlist: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	log.Infof("processed playlist:\n%+v", pl)

	err = s.wp.Post(pl)
	if err != nil {
		log.Errorf("could not post playlist: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	log.Infof("posted playlist!")

}
