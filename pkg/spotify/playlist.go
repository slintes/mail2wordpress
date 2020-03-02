package spotify

import (
	"fmt"
	"html"
	"strings"
	"time"

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

	songIDs := make([]spotify.ID, 0)
	for _, song := range pl.Songs {
		if id, err := sp.searchSong(song); err == nil {
			songIDs = append(songIDs, id)
		} else {
			fmt.Printf("error finding song: %v", err)
		}
		time.Sleep(200 * time.Millisecond)
	}

	if _, err = sp.client.AddTracksToPlaylist(newPl.ID, songIDs...); err != nil {
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
		return "", fmt.Errorf("song not found: %v", song.Title)
	}

	return id, nil
}
