package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/slintes/bluesstammtisch/pkg/playlist"
	"github.com/slintes/bluesstammtisch/pkg/server"
	"github.com/slintes/bluesstammtisch/pkg/spotify"
	"github.com/slintes/bluesstammtisch/pkg/wordpress"
)

var (
	debug bool
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")

	rootCmd.AddCommand(wordpressCmd)
	rootCmd.AddCommand(spotifyCmd)
	rootCmd.AddCommand(localPlCmd)

	//log.SetFormatter(&log.JSONFormatter{})
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}

var rootCmd = &cobra.Command{
	Use:   "bluesstammtisch",
	Short: "bluesstammtisch creates wordpress or spotify playlists",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("you need to add a command: wordpress or spotify")
		return nil
	},
}

var wordpressCmd = &cobra.Command{
	Use:   "wordpress",
	Short: "wordpress processes HTTP request and creates a Wordpress post",
	RunE: func(cmd *cobra.Command, args []string) error {
		if debug {
			log.SetLevel(log.DebugLevel)
		}
		plh := playlist.NewHandler()
		wp, err := wordpress.NewPoster()
		if err != nil {
			return err
		}
		server.NewServer(plh, wp).Serve()
		return nil
	},
}

var localPlCmd = &cobra.Command{
	Use:   "local-playlist",
	Short: "local-playlist processes a local playlist file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if debug {
			log.SetLevel(log.DebugLevel)
		}
		plh := playlist.NewHandler()
		pl, err := plh.Process("/home/msluiter/Downloads/Playlist.csv")
		if err != nil {
			return err
		}
		fmt.Print(pl.Body)
		return nil
	},
}

var spotifyCmd = &cobra.Command{
	Use:   "spotify",
	Short: "spotify processes file and creates a Spotify playlist",
	RunE: func(cmd *cobra.Command, args []string) error {
		if debug {
			log.SetLevel(log.DebugLevel)
		}
		plh := playlist.NewHandler()
		pl, err := plh.Process("/home/msluiter/Downloads/Playlist.csv")
		if err != nil {
			return err
		}
		sp, err := spotify.NewSpotify()
		if err != nil {
			return err
		}
		err = sp.CreatePlaylist(pl)
		if err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	if len(os.Args) == 1 {
		rootCmd.SetArgs([]string{"wordpress"})
	}
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("something went terrible wrong: %v", err)
	}
}
