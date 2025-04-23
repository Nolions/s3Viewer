package main

import (
	"github.com/Nolions/s3Viewer/config"
	"github.com/Nolions/s3Viewer/internal/tui"
)

func main() {
	conf := config.NewAWSConfig()

	s := tui.NewS3App(conf)
	s.BuildUI()

	if err := s.App.Run(); err != nil {
		panic(err)
	}
}
