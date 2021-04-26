package aws

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/golang/mock/gomock"
	"github.com/maruina/playground/mocks"
	. "github.com/onsi/gomega"
)

type FakeS3Client struct{}

func (f FakeS3Client) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	// If some parameters are missing, return an error
	if aws.ToString(params.Bucket) == "" {
		return nil, fmt.Errorf("target bucket is missing")
	}
	if aws.ToString(params.Key) == "" {
		return nil, fmt.Errorf("bucket key is missing")
	}
	return &s3.GetObjectOutput{
		Body: ioutil.NopCloser(bytes.NewReader([]byte("this is the body foo bar baz"))),
	}, nil
}
func TestGetObjectFromS3(t *testing.T) {
	testCases := []struct {
		name          string
		bucket        string
		key           string
		expectedError error
		expected      []byte
	}{
		{
			name:          "happy path",
			bucket:        "foo",
			key:           "bar",
			expectedError: nil,
			expected:      []byte("this is the body foo bar baz"),
		},
		{
			name:          "bucket is missing",
			bucket:        "",
			key:           "bar",
			expectedError: errors.New("target bucket is missing"),
			expected:      nil,
		},
		{
			name:          "key is missing",
			bucket:        "foo",
			key:           "",
			expectedError: errors.New("bucket key is missing"),
			expected:      nil,
		},
	}
	ctx := context.TODO()
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := GetS3Object(ctx, FakeS3Client{}, testCase.bucket, testCase.key)
			// Needed to use gomega with the testing package
			g := NewWithT(t)
			g.Expect(got).To(Equal(testCase.expected))
			// In gomega you can't use Equal for something that is nil,
			// this is why we need the if - else clause depending on the err
			if err == nil {
				g.Expect(err).To(BeNil())
			} else {
				g.Expect(err).To(Equal(testCase.expectedError))
			}
		})
	}
}

func TestGetObjectFromS3HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mocks.NewMockS3Client(ctrl)

	ctx := context.Background()
	bucket := "bucket"
	key := "key"
	expectedBody := []byte("this is the body")
	fakeBucket := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	fakeObject := &s3.GetObjectOutput{
		Body: ioutil.NopCloser(bytes.NewReader(expectedBody)),
	}

	m.EXPECT().GetObject(ctx, fakeBucket).Return(fakeObject, nil)
	g := NewWithT(t)
	got, err := GetS3Object(ctx, m, bucket, key)
	g.Expect(got).To(Equal(expectedBody))
	g.Expect(err).To(BeNil())

}
