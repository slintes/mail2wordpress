package spotify

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	sp "github.com/zmb3/spotify"
)

const redirectURI = "http://localhost:8080/callback"

var (
	id     = os.Getenv("SPOTIFY_ID")
	secret = os.Getenv("SPOTIFY_SECRET")
	User   = os.Getenv("SPOTIFY_USERNAME")

	auth     = sp.NewAuthenticator(redirectURI, sp.ScopePlaylistModifyPrivate)
	chClient = make(chan *sp.Client)
	chErr    = make(chan error)
	state    = "abc123"
)

type Spotify struct {
	client *sp.Client
	user   *sp.PrivateUser
}

func NewSpotify() (*Spotify, error) {
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	var client *sp.Client
	select {
	case client = <-chClient:
	case err := <-chErr:
		return nil, err
	case <-time.After(20 * time.Second):
		return nil, fmt.Errorf("login timed out")
	}

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		return nil, err
	}
	fmt.Println("You are logged in as:", user.ID)

	spotify := Spotify{
		client: client,
		user:   user,
	}

	return &spotify, nil
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		chErr <- err
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		chErr <- fmt.Errorf("wrong state field")
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	chClient <- &client
}
