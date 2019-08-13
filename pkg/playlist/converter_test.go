package playlist

import (
	"github.com/slintes/bluesstammtisch/pkg/types"
	"reflect"
	"testing"
)

func TestHandler_convert(t *testing.T) {

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
				csv: `Playlist 26.06.19 Bluesstammtisch Ems Vechte Welle (77);;;
;;;
Interpret;Album;Titel;Lable
Tedeschi Trucks Band;Shine;Shame;Concord
;;Shine, Hightime;
;;;
` + ";;;\r\n;;;\r\n", // test windows line feed, should not result in empty lines
			},
			want: &types.Playlist{
				Title: "Playlist 26.06.2019",
				Body: `
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
<h2>Die nächste Sendung ist am <strong>10.07.</strong>!</h2>
&nbsp;
`,
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
				t.Errorf("convert() got = %v, want %v", got, tt.want)
			}
		})
	}
}
