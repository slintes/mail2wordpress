package playlist

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/slintes/bluesstammtisch/pkg/types"
)

func TestHandler_convert(t *testing.T) {

	playlistDate := "26.06.19"
	expectedDate, _ := time.Parse("02.01.06", playlistDate)

	type args struct {
		csv string
	}
	tests := []struct {
		name    string
		args    args
		want    *types.Playlist
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				csv: `Playlist ` + playlistDate + ` Bluesstammtisch Ems Vechte Welle (77);;;
;;;
Moderation:;Gerd;N�chster Bluesstammtisch:;22.04.2020
Interpret;Album;Titel;Lable
Tedeschi Trucks Band;Shine;Shame;Concord
;;Shine, Hightime;
;;;
` + ";;;\r\n;;;\r\n", // test windows line feed, should not result in empty lines
			},
			want: &types.Playlist{
				Title: "Playlist 26.06.2019",
				Date:  expectedDate,
				Body: `
<h2><i>Moderation: <strong>Gerd</strong></i></h2>
&nbsp;<br>
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
    <tr><td>Tedeschi Trucks Band</td><td>Shine</td><td>Shame</td><td>Concord</td></tr>
    <tr><td>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"</td><td>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"</td><td>Shine, Hightime</td><td>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"</td></tr>
  </tbody>
</table>
<h2>Die nächste Sendung ist am <strong>22.04.2020</strong>!</h2>
&nbsp;<br>
`,
				Songs: []types.Song{
					{
						Artist: "Tedeschi Trucks Band",
						Album:  "Shine",
						Title:  "Shame",
						Label:  "Concord",
					},
					{
						Artist: "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;\"",
						Album:  "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;\"",
						Title:  "Shine, Hightime",
						Label:  "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;\"",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{}
			got, err := h.convert(tt.args.csv)
			if (err != nil) != tt.wantErr {
				t.Errorf("convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convert() got:\n%+v\n\n, want:\n%+v", got, tt.want)
			}
		})
	}
}

func TestHandler_findDate(t *testing.T) {
	s := []string{"Playlist 02.03.04 Bluesstammtisch xyz", ""}

	h := &Handler{}
	found, date := h.findDate(s)
	if !found {
		t.Errorf("date not found")
	}
	if date == nil {
		t.Errorf("date is nil?!")
	}
	if date.Year() != 2004 || date.Month() != 03 || date.Day() != 2 {
		t.Errorf("wrong date parsed: %v", date)
	}

	fmt.Printf("found date (formatted now): %s\n", date.Format("02.01.2006"))
}
