package spotify

import (
	"fmt"
	"html"
	"strings"

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

	// remove some chars and words which might not be exactly the same on spotify
	artist = strings.Replace(artist, " and ", " ", -1)
	artist = strings.Replace(artist, "& ", " ", -1)
	artist = strings.Replace(artist, ", ", " ", -1)
	artist = strings.Replace(artist, " ft ", " ", -1)
	artist = strings.Replace(artist, " feat ", " ", -1)
	artist = strings.Replace(artist, " Feat ", " ", -1)
	artist = strings.Replace(artist, " ft. ", " ", -1)
	artist = strings.Replace(artist, " feat. ", " ", -1)
	artist = strings.Replace(artist, " Feat. ", " ", -1)
	artist = strings.Replace(artist, " featuring ", " ", -1)
	artist = strings.Replace(artist, " Featuring ", " ", -1)

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
	for _, track := range result.Tracks.Tracks {

		// check if we have the right track in case we found many tracks from an album with the same name
		if strings.ToLower(track.Name[0:3]) != strings.ToLower(title[0:3]) {
			continue
		}

		// remember 1st as best match if album not found
		if id == "" {
			id = track.ID
		}

		fmt.Printf("    album: %v\n", track.Album.Name)

		if strings.ToLower(track.Album.Name) == strings.ToLower(album) {
			return track.ID, nil
		}
	}

	if string(id) == "" {
		return "", fmt.Errorf("no song found: %+v", song)
	}

	return id, nil
}
