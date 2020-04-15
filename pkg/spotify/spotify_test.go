package spotify

import (
	"fmt"
	"github.com/zmb3/spotify"
	"log"
	"testing"
)

func XTestNewSpotify(t *testing.T) {
	sp, err := NewSpotify()
	if err != nil {
		log.Fatalf("couldn't get client: %v", err)
	}

	playlistPage, err := sp.client.CurrentUsersPlaylists()
	if err != nil {
		log.Fatalf("couldn't get user's playlists: %v", err)
	}

	for _, playlist := range playlistPage.Playlists {
		fmt.Println("  ", playlist.Name)

		fullPlaylist, err := sp.client.GetPlaylist(playlist.ID)
		if err != nil {
			log.Fatalf("couldn't get playlist: %v", err)
		}
		for _, track := range fullPlaylist.Tracks.Tracks {
			fmt.Println("    ", track.Track.Name)
		}

	}

	newPl, err := sp.client.CreatePlaylistForUser(sp.user.ID, "Playlist 11.22", "Bluesstammtisch am 33.44", false)
	if err != nil {
		log.Fatalf("couldn't create playlists: %v", err)
	}
	newPlID := newPl.ID

	result, err := sp.client.Search(`way back home artist:Gregor Hilden album:blue in red`, spotify.SearchTypeTrack)
	if err != nil {
		log.Fatalf("couldn't search track: %v", err)
	}

	for i, track := range result.Tracks.Tracks {
		fmt.Println("     " + track.Name)
		if i == 0 {
			_, err = sp.client.AddTracksToPlaylist(newPlID, track.ID)
			if err != nil {
				log.Fatalf("couldn't add track: %v", err)
			}
		}
	}

}
