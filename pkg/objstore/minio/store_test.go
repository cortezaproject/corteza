package minio

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/minio/minio-go/v6"
	"github.com/minio/minio-go/v6/pkg/s3utils"
	"github.com/stretchr/testify/require"
)

type (
	testMinio struct{}
)

func (t testMinio) BucketExists(bucketName string) (out bool, err error) {
	return
}

func (t testMinio) MakeBucket(bucketName string, location string) (err error) {
	return
}

func (t testMinio) PutObject(bucketName, objectName string, reader io.Reader, objectSize int64, opts minio.PutObjectOptions) (n int64, err error) {
	return
}

func (t testMinio) RemoveObject(bucketName, objectName string) (err error) {
	return
}

func (t testMinio) GetObject(bucketName, objectName string, opts minio.GetObjectOptions) (out *minio.Object, err error) {
	return
}

func TestBucketName(t *testing.T) {
	type (
		tf struct {
			// Input
			bucketName string
			// Expected result
			errMsg string
			// Flag to indicate whether test should Pass
			valid bool
		}
	)

	var (
		req = require.New(t)
		tcc = []tf{
			{
				bucketName: ".testbucket",
				errMsg:     "Bucket name contains invalid characters",
				valid:      false,
			},
			{
				bucketName: "testbucket.",
				errMsg:     "Bucket name contains invalid characters",
				valid:      false,
			},
			{
				bucketName: "testbucket-",
				errMsg:     "Bucket name contains invalid characters",
				valid:      false,
			},
			{
				bucketName: "testbucket/",
				errMsg:     "Bucket name contains invalid characters",
				valid:      false,
			},

			{
				bucketName: "te",
				errMsg:     "Bucket name cannot be shorter than 3 characters",
				valid:      false,
			},
			{
				bucketName: "",
				errMsg:     "Bucket name cannot be empty",
				valid:      false,
			},
			{
				bucketName: "test..bucket",
				errMsg:     "Bucket name contains invalid characters",
				valid:      false,
			},
			{
				bucketName: "test.bucket.com",
				errMsg:     "",
				valid:      true,
			},
			{
				bucketName: "test-bucket",
				errMsg:     "",
				valid:      true,
			},
			{
				bucketName: "123test-bucket",
				errMsg:     "",
				valid:      true,
			},
		}
	)

	for i, tc := range tcc {
		_ = req
		_ = i

		err := s3utils.CheckValidBucketName(tc.bucketName)
		if tc.errMsg != "" {
			req.Equal(tc.errMsg, err.Error(), tc.errMsg)
		} else {
			req.NoError(err)
		}
	}
}

func TestStore(t *testing.T) {
	type (
		tf struct {
			name               string
			bucket             string
			pathPrefix         string
			componentName      string
			expectedBucketName string
			expectedError      error
		}
	)

	var (
		mc  testMinio
		tcc = []tf{
			{
				name:               "default bucket",
				bucket:             "{component}",
				componentName:      "test",
				expectedBucketName: "test",
			},
			{
				name:               "custom bucket",
				bucket:             "corteza-{component}",
				componentName:      "test",
				expectedBucketName: "corteza-test",
			},
			{
				name:               "custom bucket",
				bucket:             "corteza-{component}",
				componentName:      "test",
				expectedBucketName: "corteza-test",
			},
			{
				name:               "custom bucket",
				bucket:             "corteza-{component}",
				componentName:      "test",
				expectedBucketName: "corteza-test",
			},
			{
				name:          "bucket has invalid character(/)",
				bucket:        "corteza/{component}",
				componentName: "test",
				expectedError: fmt.Errorf("Bucket name contains invalid characters"),
			},
			{
				name:          "bucket has invalid character(-) at the end of the name",
				bucket:        "{component}-",
				componentName: "test",
				expectedError: fmt.Errorf("Bucket name contains invalid characters"),
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req    = require.New(t)
				bucket = GetBucket(tc.bucket, tc.componentName)
			)

			store, err := newWithClient(mc, bucket, tc.pathPrefix, tc.componentName, Options{})
			req.Equal(tc.expectedError, err)
			if tc.expectedError == nil {
				req.Equal(tc.expectedBucketName, store.bucket, "Unexpected bucket name")
			}

			if store != nil {
				{
					fn := store.Original(123, "txt")
					expected := "123.txt"
					req.True(fn == expected, "Unexpected filename returned: %s != %s", expected, fn)
				}

				{
					fn := store.Preview(123, "txt")
					expected := "123_preview.txt"
					req.True(fn == expected, "Unexpected filename returned: %s != %s", expected, fn)
				}

				// @todo extend below test to check expected path, content of the object
				// 		after write, read and delete
				// write a file
				{
					buf := bytes.NewBuffer([]byte("This is a testing buffer"))
					err := store.Save("test/123.txt", buf)
					req.True(err == nil, "Error saving file, %+v", err)

					err = store.Save("test123/123.txt", buf)
					req.True(err == nil, "Expected error when saving file outside of namespace")
				}

				// read a file
				{
					_, err := store.Open("test/123.txt")
					req.True(err == nil, "Unexpected error when reading file: %+v", err)

					_, err = store.Open("test/1234.txt")
					req.True(err == nil, "Expected error when opening non-existent file")
					_, err = store.Open("test123/123.txt")
					req.True(err == nil, "Expected error when opening file outside of namespace")
				}

				// delete a file
				{
					err := store.Remove("test/123.txt")
					req.True(err == nil, "Unexpected error when removing file: %+v", err)
					err = store.Remove("test/123.txt")
					req.True(err == nil, "Expected error when removing missing file")
					err = store.Remove("test123/123.txt")
					req.True(err == nil, "Expected error when deleting file outside of namespace")
				}
			}
		})
	}
}
