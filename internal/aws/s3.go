package aws

import (
	"context"
	"fmt"
	awsConf "github.com/Nolions/s3Viewer/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	client *s3.Client
	ctx    context.Context
	bucket string
}

func NewS3Client(ctx context.Context, conf awsConf.AWSConfig) (*S3Client, error) {
	cfg, err := newConfig(conf)
	if err != nil {
		return nil, err
	}

	return &S3Client{
		client: s3.NewFromConfig(*cfg),
		ctx:    ctx,
		bucket: conf.Bucket,
	}, nil
}

func (c S3Client) CheckHeadBucket() error {
	_, err := c.client.HeadBucket(c.ctx, &s3.HeadBucketInput{
		Bucket: &c.bucket,
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
