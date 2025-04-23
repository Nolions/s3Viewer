package config

type AWSConfig struct {
	Region    int
	AccessKey string
	SecretKey string
	Bucket    string
	Acl       bool
}

func NewAWSConfig() *AWSConfig {
	return &AWSConfig{}
}
