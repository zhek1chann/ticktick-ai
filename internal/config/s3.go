package config

import "os"

type S3Config struct {
	endpoint  string
	region    string
	accessKey string
	secretKey string
	bucket    string
}

func (c S3Config) Endpoint() string  { return c.endpoint }
func (c S3Config) Region() string    { return c.region }
func (c S3Config) AccessKey() string { return c.accessKey }
func (c S3Config) SecretKey() string { return c.secretKey }
func (c S3Config) Bucket() string    { return c.bucket }

func newS3ConfigEnv() (S3Config, error) {
	return S3Config{
		endpoint:  os.Getenv("S3_ENDPOINT"),
		region:    os.Getenv("S3_REGION"),
		accessKey: os.Getenv("S3_ACCESS_KEY"),
		secretKey: os.Getenv("S3_SECRET_KEY"),
		bucket:    os.Getenv("S3_BUCKET"),
	}, nil
}
