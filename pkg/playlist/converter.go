package playlist

import (
	"fmt"
	"html"
	"html/template"
	"strings"
	"time"

	"github.com/slintes/bluesstammtisch/pkg/types"
)

const (
	Ditto = template.HTML("&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;\"")
)

type templateData struct {
	Songs []types.Song
	Next  string
}

func (h *Handler) convert(csv string) (*types.Playlist, error) {

	songs, err := h.parseCsv(csv)
	if err != nil {
		return nil, fmt.Errorf("could not parse csv: %v", err)
	}

	data := &templateData{
		Songs: songs,
		Next:  h.getNextDate(),
	}

	t := template.Must(template.New("playlist").Parse(h.getTemplate()))
	buf := strings.Builder{}
	if err = t.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("error executing template: %v", err)
	}

	return &types.Playlist{
		Title: fmt.Sprintf("Playlist %s", h.getLastDate()),
		Body:  buf.String(),
		Songs: songs,
	}, nil
}

func (h *Handler) parseCsv(csv string) ([]types.Song, error) {

	var songs []types.Song

	lines := strings.Split(csv, "\n")

	lastArtist := template.HTML("")
	lastAlbum := template.HTML("")
	lastLabel := template.HTML("")

	for _, line := range lines {

		// there is someone sending the playlist with different columns...
		// TODO detect this
		fromWinnie := false

		minLength := 3
		colArtist := 0
		colAlbum := 1
		colTitle := 2
		colLabel := 3

		if fromWinnie {
			minLength = 4
			colTitle = 3
			colLabel = 4
		}

		parts := strings.Split(line, ";")
		if len(parts) < minLength {
			continue
		}

		if h.skipLine(parts) {
			continue
		}

		newArtist := false
		artist := template.HTML(html.EscapeString(parts[colArtist]))
		if h.isDitto(newArtist, lastArtist, artist) {
			artist = Ditto
		} else {
			lastArtist = artist
			newArtist = true
		}

		album := template.HTML(html.EscapeString(parts[colAlbum]))
		if h.isDitto(newArtist, lastAlbum, album) {
			album = Ditto
		} else {
			lastAlbum = album
		}

		title := template.HTML(html.EscapeString(parts[colTitle]))

		label := template.HTML("")
		if len(parts) >= colLabel+1 {
			label = template.HTML(html.EscapeString(parts[colLabel]))
		}
		if h.isDitto(newArtist, lastLabel, label) {
			label = Ditto
		} else {
			lastLabel = label
		}

		songs = append(songs, types.Song{
			Artist: artist,
			Album:  album,
			Title:  title,
			Label:  label,
		})

	}

	return songs, nil
}

func (h *Handler) skipLine(parts []string) bool {
	// skip lines without songs
	if strings.HasPrefix(parts[0], "Playlist") {
		return true
	}
	if strings.Contains(parts[0], "Bluesstammtisch") {
		return true
	}
	if strings.HasPrefix(parts[0], "Interpret") {
		return true
	}
	if len(strings.TrimSpace(parts[0])) == 0 &&
		len(strings.TrimSpace(parts[1])) == 0 &&
		len(strings.TrimSpace(parts[2])) == 0 &&
		len(strings.TrimSpace(parts[3])) == 0 {

		return true

	}
	return false
}

func (h *Handler) isDitto(newArtist bool, oldVal template.HTML, newVal template.HTML) bool {
	if !newArtist &&
		(string(newVal) == string(oldVal) ||
			strings.TrimSpace(string(newVal)) == "" ||
			strings.TrimSpace(string(newVal)) == "\"" ||
			strings.HasPrefix(string(newVal), "\"\"")) {
		return true
	}
	return false
}

func (h *Handler) getLastDate() string {
	return h.getLastWednesday().Format("02.01.2006")
}

func (h *Handler) getNextDate() string {
	// get wednesday in 2 weeks
	lastWed := h.getLastWednesday()
	nextWed := lastWed.AddDate(0, 0, 14)
	return nextWed.Format("02.01.")
}

func (h *Handler) getLastWednesday() time.Time {
	// get last wednesday (today is fine)
	now := time.Now()
	day := int(now.Weekday())
	// if sunday to tuesday, add a week
	if day < 3 {
		day += 7
	}
	lastWed := now.AddDate(0, 0, -(day - 3))
	return lastWed
}

func (h *Handler) getTemplate() string {
	return `
<table class="table table-responsive">
  <thead>
    <tr>
      <th>Interpret</th>
      <th>Album</th>
      <th>Titel</th>
      <th>Label</th>
    </tr>
  </thead>
  <tbody>{{range .Songs}}
    <tr><td>{{.Artist}}</td><td>{{.Album}}</td><td>{{.Title}}</td><td>{{.Label}}</td></tr>{{end}}
  </tbody>
</table>
<h2>Die nächste Sendung ist am <strong>{{.Next}}</strong>!</h2>
&nbsp;
`
}
