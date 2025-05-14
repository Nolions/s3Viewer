package main

import (
	"context"
	"github.com/Nolions/s3Viewer/config"
	"github.com/Nolions/s3Viewer/internal/tui"
)

func main() {
	conf := config.NewAWSConfig()
	ctx := context.Background()
	s := tui.NewS3App(ctx, conf)
	s.BuildUI()

	if err := s.App.Run(); err != nil {
		panic(err)
	}
}
