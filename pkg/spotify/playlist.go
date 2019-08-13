package spotify

import (
	"fmt"
	"html"

	"github.com/zmb3/spotify"

	"github.com/slintes/bluesstammtisch/pkg/playlist"
	"github.com/slintes/bluesstammtisch/pkg/types"
)

var (
	lastArtist string
	lastAlbum  string
)

func (sp *Spotify) CreatePlaylist(pl *types.Playlist) error {

	newPl, err := sp.client.CreatePlaylistForUser(sp.user.ID, pl.Title, "", false)
	if err != nil {
		return err
	}

	for _, song := range pl.Songs {
		if sp.addSong(newPl.ID, song); err != nil {
			return err
		}
	}

	return nil
}

func (sp *Spotify) addSong(plID spotify.ID, song types.Song) error {
	id, err := sp.searchSong(song)
	if err != nil {
		return err
	}
	if _, err = sp.client.AddTracksToPlaylist(plID, id); err != nil {
		return err
	}
	return nil
}

func (sp *Spotify) searchSong(song types.Song) (spotify.ID, error) {

	title := html.UnescapeString(string(song.Title))
	artist := html.UnescapeString(string(song.Artist))
	album := html.UnescapeString(string(song.Album))

	if song.Artist == playlist.Ditto {
		artist = lastArtist
	} else {
		lastArtist = artist
	}
	if song.Album == playlist.Ditto {
		album = lastAlbum
	} else {
		lastAlbum = album
	}

	fmt.Printf("searching %v from %v on %v", title, artist, album)

	result, err := sp.client.Search(fmt.Sprintf("%s artist:%s\n", title, artist), spotify.SearchTypeTrack)
	if err != nil {
		return "", err
	}

	fmt.Printf("  found %v tracks\n", len(result.Tracks.Tracks))

	var id spotify.ID
	for i, track := range result.Tracks.Tracks {
		// remember 1st as best match if album not found
		if i == 0 {
			id = track.ID
		}

		fmt.Printf("    album: %v\n", track.Album.Name)

		if track.Album.Name == album {
			return track.ID, nil
		}
	}

	if string(id) == "" {
		return "", fmt.Errorf("no song found: %+v", song)
	}

	return id, nil
}
