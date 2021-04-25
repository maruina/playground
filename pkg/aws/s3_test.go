package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/golang/mock/gomock"
	"github.com/maruina/playground/mocks"
)

func TestS3Client(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	mockS3 := mocks.NewMockS3Client(mockCtrl)

	wantBucket := &s3.GetObjectInput{
		Bucket: aws.String("bucket"),
		Key:    aws.String("key"),
	}

	mockS3.EXPECT().GetObject(context.Background(), wantBucket).Return([]byte("this is the body foo bar baz"), nil)

}
