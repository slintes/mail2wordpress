package wordpress

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	wp "github.com/robbiet480/go-wordpress"
	log "github.com/sirupsen/logrus"
	"github.com/slintes/bluesstammtisch/pkg/types"
)

var (
	wpUrl    = os.Getenv("WP_URL")
	wpUser   = os.Getenv("WP_USER")
	wpPwd    = os.Getenv("WP_PWD")
	wpCat, _ = strconv.Atoi(os.Getenv("WP_CAT"))
	wpAuthor int
)

type Poster struct {
	client *wp.Client
}

func NewPoster() (*Poster, error) {

	// create wp-api client
	tp := wp.BasicAuthTransport{
		Username: wpUser,
		Password: wpPwd,
	}

	client, err := wp.NewClient(wpUrl, tp.Client())
	if err != nil {
		return nil, fmt.Errorf("could not connect: %v", err)
	}

	// get user id
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	currentUser, resp, err := client.Users.Me(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	if resp != nil && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status: %v %v", resp.StatusCode, resp.Status)
	}
	wpAuthor = currentUser.ID

	return &Poster{
		client,
	}, nil

}

func (p *Poster) Post(pl *types.Playlist) error {
	post := &wp.Post{
		Title: wp.RenderedString{
			Raw: pl.Title,
		},
		Content: wp.RenderedString{
			Raw: pl.Body,
		},
		Categories: []int{wpCat},
		Format:     wp.PostFormatStandard,
		Type:       wp.PostTypePost,
		Status:     wp.PostStatusDraft,
		Author:     wpAuthor,
		Date:       wp.Time{Time: time.Now().UTC()},
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	newPost, resp, err := p.client.Posts.Create(ctx, post)
	if err != nil {
		return fmt.Errorf("error posting playlist: %v", err)
	}
	if resp != nil && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("wrong status: %v %v", resp.StatusCode, resp.Status)
	}
	log.Infof("new post ID: %v", newPost.ID)
	return nil
}
