package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/slintes/mail2wordpress/pkg/playlist"
	"github.com/slintes/mail2wordpress/pkg/server"
	"github.com/slintes/mail2wordpress/pkg/wordpress"
	"github.com/spf13/cobra"
)

var (
	debug bool
)

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")

	//log.SetFormatter(&log.JSONFormatter{})
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}

var rootCmd = &cobra.Command{
	Use:   "request2wordpress",
	Short: "request2wordpress processes HTTP request and creates a Wordpress post",
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

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("something went terrible wrong: %v", err)
	}
}
