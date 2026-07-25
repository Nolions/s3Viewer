package config

type AWSConfig struct {
	Type         string // "AWS S3" or "MinIO"
	Endpoint     string // Host/Endpoint URL for MinIO
	UsePathStyle bool   // MinIO requires path style
	Region       string
	AccessKey    string
	SecretKey    string
	Bucket       string
	Acl          bool
}

func NewAWSConfig() *AWSConfig {
	return &AWSConfig{
		Type:         "AWS S3",
		Region:       "us-east-1",
		UsePathStyle: false,
	}
}
