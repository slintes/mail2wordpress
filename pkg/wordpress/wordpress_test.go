package wordpress

import (
	"fmt"
	"testing"
	"time"
)

func Test_getPublishDate(t *testing.T) {

	loc, _ := time.LoadLocation("Europe/Berlin")
	now := time.Now().In(loc)
	dateNow := now.Format("02.01.2006")
	fmt.Printf("now as string: %v\n", dateNow)

	// Test that playlist is for yesterday, publish date now
	fmt.Println("testing yesterday's playlist")
	playlistDate, _ := time.Parse("02.01.2006", dateNow)
	playlistDate = playlistDate.Add(-24 * time.Hour)
	playlistDateLoc := playlistDate.In(loc)
	fmt.Printf("local playlist date: %v\n", playlistDateLoc)

	publishDate := getPublishDate(playlistDate)
	fmt.Printf("local publish date: %v\n", publishDate)

	if publishDate.Day() != now.Day() || publishDate.Hour() != now.Hour() {
		t.Errorf("wrong publish date for yesterday's playlist: %v\n", publishDate.Hour())
	}

	// Test that playlist is for today, publish date is today 22h
	fmt.Println("testing today's playlist")
	playlistDate, _ = time.Parse("02.01.2006", dateNow)
	playlistDateLoc = playlistDate.In(loc)
	fmt.Printf("local playlist date: %v\n", playlistDateLoc)

	publishDate = getPublishDate(playlistDate)
	fmt.Printf("local publish date: %v\n", publishDate)

	if publishDate.Hour() != 22 {
		t.Errorf("wrong hour for today's playlist: %v\n", publishDate.Hour())
	}

	// Test that playlist is for tomorrow, publish date is tomorrow 22h
	fmt.Println("testing tomorrow's playlist")
	playlistDate, _ = time.Parse("02.01.2006", dateNow)
	playlistDate = playlistDate.Add(24 * time.Hour)
	playlistDateLoc = playlistDate.In(loc)
	fmt.Printf("local playlist date: %v\n", playlistDateLoc)

	publishDate = getPublishDate(playlistDate)
	fmt.Printf("local publish date: %v\n", publishDate)

	if publishDate.Day() != (now.Day()+1) || publishDate.Hour() != 22 {
		t.Errorf("wrong publish date for tomorrow's playlist: %v\n", publishDate.Hour())
	}
}
