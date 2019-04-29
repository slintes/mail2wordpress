package playlist

import (
	"fmt"
	"html"
	"html/template"
	"strings"
	"time"

	"github.com/slintes/mail2wordpress/pkg/types"
)

const (
	ditto = template.HTML("&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;\"")
)

type templateData struct {
	Songs []song
	Next  string
}

type song struct {
	Artist template.HTML
	Album  template.HTML
	Title  template.HTML
	Label  template.HTML
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
	}, nil
}

func (h *Handler) parseCsv(csv string) ([]song, error) {

	var songs []song

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

		// skip lines without songs
		if strings.HasPrefix(parts[0], "Playlist") {
			continue
		}
		if strings.Contains(parts[0], "Bluesstammtisch") {
			continue
		}
		if strings.HasPrefix(parts[0], "Interpret") {
			continue
		}
		if len(parts[0]) == 0 && len(parts[1]) == 0 && len(parts[2]) == 0 && len(parts[3]) == 0 {
			continue
		}

		newArtist := false
		artist := template.HTML(html.EscapeString(parts[colArtist]))
		if h.isDitto(newArtist, lastArtist, artist) {
			artist = ditto
		} else {
			lastArtist = artist
			newArtist = true
		}

		album := template.HTML(html.EscapeString(parts[colAlbum]))
		if h.isDitto(newArtist, lastAlbum, album) {
			album = ditto
		} else {
			lastAlbum = album
		}

		title := template.HTML(html.EscapeString(parts[colTitle]))

		label := template.HTML("")
		if len(parts) >= colLabel+1 {
			label = template.HTML(html.EscapeString(parts[colLabel]))
		}
		if h.isDitto(newArtist, lastLabel, label) {
			label = ditto
		} else {
			lastLabel = label
		}

		songs = append(songs, song{
			Artist: artist,
			Album:  album,
			Title:  title,
			Label:  label,
		})

	}

	return songs, nil
}

func (h *Handler) isDitto(newArtist bool, oldVal template.HTML, newVal template.HTML) bool {
	if !newArtist &&
		(string(newVal) == string(oldVal) ||
			strings.TrimSpace(string(newVal)) == "" ||
			string(newVal) == "\"" ||
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
  <tbody>
    {{range .Songs}}
    <tr><td>{{.Artist}}</td><td>{{.Album}}</td><td>{{.Title}}</td><td>{{.Label}}</td></tr>
	{{end}}
  </tbody>
</table>
<h2>Die nächste Sendung ist am <strong>{{.Next}}</strong>!</h2>
&nbsp;
`
}
