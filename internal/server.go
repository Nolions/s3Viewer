package internal

import "github.com/aws/aws-sdk-go-v2/service/s3"

type App struct {
	s3Client s3.Client
	bucket   string
}
