package aws

import (
	"context"
	"io/ioutil"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

// GetS3Object returns an object from an S3 bucket
func GetS3Object(ctx context.Context, client S3Client, bucket, key string) ([]byte, error) {
	object, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		bErr := object.Body.Close()
		if err == nil {
			err = bErr
		}
	}()
	return ioutil.ReadAll(object.Body)
}
