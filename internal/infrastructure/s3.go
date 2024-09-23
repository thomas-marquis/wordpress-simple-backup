package infrastructure

type S3Impl struct {
	accessKey string
	secretKey string
	region    string
	bucket    string
}

func NewS3Impl(accessKey, secretKey, region, bucket string) *S3Impl {
	return &S3Impl{
		accessKey: accessKey,
		secretKey: secretKey,
		region:    region,
		bucket:    bucket,
	}
}
