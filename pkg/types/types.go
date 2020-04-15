package types

import (
	"html/template"
	"time"
)

type Song struct {
	Artist template.HTML
	Album  template.HTML
	Title  template.HTML
	Label  template.HTML
}

type Playlist struct {
	Title string
	Date  time.Time
	Body  string
	Songs []Song
}
