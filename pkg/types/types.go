package types

import "html/template"

type Song struct {
	Artist template.HTML
	Album  template.HTML
	Title  template.HTML
	Label  template.HTML
}

type Playlist struct {
	Title string
	Body  string
	Songs []Song
}
